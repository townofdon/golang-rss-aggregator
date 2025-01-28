package scraper

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Error(fmt.Sprintf("[scraper] unable to mark feed %v as fetched - ", feed.ID) + err.Error())
		return
	}
	rssFeed, err := util.UrlToFeed(feed.Url)
	if err != nil {
		log.Error(fmt.Sprintf("[scraper] could not parse RSS feed for url %v - ", feed.Url) + err.Error())
		return
	}
	log.Info(fmt.Sprintf("[scraper] feed \"%s\" collected, %v posts found", feed.Name, len(rssFeed.Channel.Item)))
	for i, item := range rssFeed.Channel.Item {
		log.Info(fmt.Sprintf("[scraper] found post %v: %v", i+1, item.Title))
	}
}
