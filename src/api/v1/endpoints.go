package v1

import "github.com/townofdon/tutorial-go-rss-server/src/api"

type Endpoint api.Clients

func SetupEndpoints(clients *api.Clients) Endpoint {
	return Endpoint{
		DB: clients.DB,
	}
}
