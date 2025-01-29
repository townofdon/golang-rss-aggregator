package v1

import (
	"fmt"
	"net/http"

	"github.com/townofdon/tutorial-go-rss-server/internal/database"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

func (e *Endpoint) GetPostsForCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(api.CtxKeyUser{}).(database.User)

	if !ok {
		util.RespondWithError(w, 403, "no logged-in user found")
		return
	}

	db, _ := api.GetDBClient()

	limit, offset := util.GetLimitOffsetFromUrlQuery(r)

	posts, err := db.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		util.RespondWithError(w, 500, fmt.Sprintf("Could not get posts: %v", err))
		return
	}

	util.RespondWithJSON(w, 200, formatPostsForResponse(posts))
}

func formatPostsForResponse(posts []database.GetPostsForUserRow) []WirePost {
	if posts == nil {
		return make([]WirePost, 0)
	}
	output := make([]WirePost, 0, len(posts))
	for _, post := range posts {
		wirePost := WirePost{
			DatabasePost: DatabasePost(post),
			Description:  &post.Description.String,
		}
		output = append(output, wirePost)
	}
	return output
}

type DatabasePost database.GetPostsForUserRow
type WirePost struct {
	// embed all fields from database.GetPostsForUserRow
	DatabasePost
	// override `description` - JSON marshalling supports string pointers
	Description *string `json:"description"`
}
