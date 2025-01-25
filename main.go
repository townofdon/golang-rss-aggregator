package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	apiV1 "github.com/townofdon/tutorial-go-rss-server/src/api/v1"
)

func main() {
	port := api.GetPort()
	// experimented with various methods of passing the db client around - 1) via global state and init, 2) via methods
	db, _ := api.GetDBClient()
	clients := api.Clients{DB: db}
	v1 := apiV1.SetupEndpoints(&clients)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// varying ways to achieve an authenticated route - 1) via closure, 2) via middleware
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", v1.HandleHealthCheck)
	v1Router.Post("/users", v1.CreateUser)
	v1Router.Get("/users/current", api.Authorized(v1.GetUserByApiKey))
	v1Router.Post("/feeds", api.Authorized(v1.CreateFeed))
	v1Router.Get("/feeds", v1.GetAllFeeds)

	v1AdminRouter := chi.NewRouter()
	v1AdminRouter.Use(api.AuthorizedMiddleware)
	v1AdminRouter.Get("/healthz", v1.HandleHealthCheck)

	router.Mount("/v1", v1Router)
	router.Mount("/v1/admin", v1AdminRouter)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on port %s...\n", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
