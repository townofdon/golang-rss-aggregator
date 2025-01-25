package auth

import (
	"fmt"
	"net/http"

	"github.com/townofdon/tutorial-go-rss-server/internal/database"

	"github.com/townofdon/tutorial-go-rss-server/src/api"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

type AuthHandlerFunc func(http.ResponseWriter, *http.Request, *api.Clients, database.User)

func Handler(api *api.Clients, handler AuthHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := ParseApiKeyFromHeader(r.Header)

		if err != nil {
			util.RespondWithError(w, 401, err.Error())
			return
		}

		user, err := api.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			util.RespondWithError(w, 404, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		handler(w, r, api, user)
	}
}
