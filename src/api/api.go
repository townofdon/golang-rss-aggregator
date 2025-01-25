package api

import (
	"net/http"

	"github.com/townofdon/tutorial-go-rss-server/internal/database"
)

type Clients struct {
	DB *database.Queries
}

type ApiHandlerFunc func(http.ResponseWriter, *http.Request, *Clients)

func Handler(api *Clients, handler ApiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, api)
	}
}
