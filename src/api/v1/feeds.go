package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/townofdon/tutorial-go-rss-server/internal/database"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

func (e *Endpoint) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramsDef struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := paramsDef{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		util.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Name == "" {
		util.RespondWithError(w, 400, "name is required")
		return
	}

	if params.Url == "" {
		util.RespondWithError(w, 400, "url is required")
		return
	}

	api.WithDBClient(func(db *database.Queries) {
		feed, err := db.CreateFeed(r.Context(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Url:       params.Url,
			UserID:    user.ID,
		})

		if err != nil {
			util.RespondWithError(w, 500, fmt.Sprintf("Error creating feed: %v", err))
			return
		}

		util.RespondWithJSON(w, 201, feed)
	})
}

func (api *Endpoint) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := api.DB.GetFeeds(r.Context())

	if err != nil {
		util.RespondWithError(w, 500, fmt.Sprintf("Error getting feeds: %v", err))
		return
	}

	util.RespondWithJSON(w, 200, feeds)
}
