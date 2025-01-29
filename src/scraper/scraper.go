package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/townofdon/tutorial-go-rss-server/internal/database"
	"github.com/townofdon/tutorial-go-rss-server/src/log"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

func Start(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Info(fmt.Sprintf("[scraper] scraping %v threads every %s", concurrency, timeBetweenRequests))
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Error("[scraper] " + err.Error())
			continue
		}
		log.Info(fmt.Sprintf("[scraper] fetching next %v rss feeds", len(feeds)))
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			// immediately increment counter to achieve waiting
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}
		// waits until wait conter is zero
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	// decrements wait counter by 1
	defer wg.Done()
	rssFeed, err := util.UrlToFeed(feed.Url)
	if err != nil {
		log.Error(fmt.Sprintf("[scraper] could not parse RSS feed for url %v - ", feed.Url) + err.Error())
		return
	}
	log.Info(fmt.Sprintf("[scraper] feed \"%s\" collected, %v posts found", feed.Name, len(rssFeed.Channel.Item)))

	_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Error(fmt.Sprintf("[scraper] unable to mark feed %v as fetched - ", feed.Name) + err.Error())
		return
	}

	numDuplicates := 0
	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := parsePostDate(item.PubDate)
		if err != nil {
			log.Error(fmt.Sprintf(
				"[scraper] could not parse PubDate \"%v\" for post \"%v\" for feed \"%v\"",
				item.PubDate,
				item.Title,
				feed.Name,
			))
		}

		post, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
				numDuplicates++
				continue
			}
			log.Error(fmt.Sprintf("[scraper] unable to save post \"%v\" for feed \"%v\" - ", item.Title, feed.Name) + err.Error())
			continue
		}
		log.Info(fmt.Sprintf("[scraper] saved new post \"%v\" (%v) for feed \"%v\"", item.Title, post.ID, feed.Name))
	}

	if numDuplicates > 0 {
		log.Info(fmt.Sprintf("[scraper] skipped %v duplicate posts for feed \"%v\"", numDuplicates, feed.Name))
	}
}

// parse a post date from an RSS feed
// see: https://pkg.go.dev/time#Layout
func parsePostDate(date string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		time.RFC3339Nano,
		// time.RFC850,
		// time.RFC822Z,
		// time.RFC822,
		// time.RubyDate,
		// time.UnixDate,
		// time.ANSIC,
	}
	var parsed time.Time
	var err error

	for _, format := range formats {
		parsed, err = time.Parse(format, date)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, err
}
