package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	v1 "github.com/townofdon/tutorial-go-rss-server/src/api/v1"
	"github.com/townofdon/tutorial-go-rss-server/src/auth"
)

func main() {
	clients, port := api.SetupApiClients()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", v1.HandleHealthCheck)
	v1Router.Post("/users", api.Handler(&clients, v1.CreateUser))
	v1Router.Get("/users/current", auth.Handler(&clients, v1.GetUserByApiKey))
	v1Router.Post("/feeds", auth.Handler(&clients, v1.CreateFeed))

	router.Mount("/v1", v1Router)

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
