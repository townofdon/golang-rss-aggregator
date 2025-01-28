package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/townofdon/tutorial-go-rss-server/internal/database"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

func (e *Endpoint) CreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	type paramsDef struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	params := paramsDef{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		util.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, ok := r.Context().Value(api.CtxKeyUser{}).(database.User)

	if !ok {
		util.RespondWithError(w, 403, "no logged-in user found")
		return
	}

	db, _ := api.GetDBClient()

	feedFollow, err := db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		statusCode := 500
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			statusCode = 400
		}
		util.RespondWithError(w, statusCode, fmt.Sprintf("Error creating feed_follow: %v", err))
		return
	}

	util.RespondWithJSON(w, 201, feedFollow)
}

func (e *Endpoint) GetFeedFollowsForCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(api.CtxKeyUser{}).(database.User)

	if !ok {
		util.RespondWithError(w, 403, "no logged-in user found")
		return
	}

	db, _ := api.GetDBClient()

	feeds, err := db.GetFeedFollowsByUserId(r.Context(), user.ID)

	if err != nil {
		util.RespondWithError(w, 500, fmt.Sprintf("Could not get feed_follows: %v", err))
		return
	}

	if feeds == nil {
		feeds = make([]database.GetFeedFollowsByUserIdRow, 0)
	}

	util.RespondWithJSON(w, 201, feeds)
}

func (e *Endpoint) DeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowId, err := uuid.Parse(chi.URLParam(r, "feedFollowId"))

	if err != nil {
		util.RespondWithError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
		return
	}

	user, ok := r.Context().Value(api.CtxKeyUser{}).(database.User)

	if !ok {
		util.RespondWithError(w, 403, "no logged-in user found")
		return
	}

	db, _ := api.GetDBClient()

	err = db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		util.RespondWithError(w, 500, fmt.Sprintf("Unable to delete feed_follow: %v", err))
		return
	}

	util.RespondWithJSON(w, 200, struct{}{})
}
