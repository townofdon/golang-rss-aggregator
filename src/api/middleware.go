package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/townofdon/tutorial-go-rss-server/internal/database"

	"github.com/townofdon/tutorial-go-rss-server/src/auth"
	"github.com/townofdon/tutorial-go-rss-server/src/util"
)

type AuthHandlerFunc func(http.ResponseWriter, *http.Request, database.User)

type CtxKeyUser struct{}

func Authorized(handler AuthHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.ParseApiKeyFromHeader(r.Header)

		if err != nil {
			util.RespondWithError(w, 401, err.Error())
			return
		}

		db, _ := GetDBClient()

		user, err := db.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			util.RespondWithError(w, 404, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyUser{}, user)
		handler(w, r.WithContext(ctx), user)
	}
}

func AuthorizedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.ParseApiKeyFromHeader(r.Header)

		if err != nil {
			util.RespondWithError(w, 401, err.Error())
			return
		}

		db, _ := GetDBClient()
		user, err := db.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			util.RespondWithError(w, 403, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyUser{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
