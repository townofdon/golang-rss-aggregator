package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	v1 "github.com/townofdon/tutorial-go-rss-server/v1"
)

func main() {
	if err := godotenv.Load(".env"); err != nil { panic(err.Error()) }

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var required")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", v1.HandleHealthCheck)
	// v1Router.Post("/healthz", v1.HandleHealthCheck)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

	fmt.Printf("Server running on port %s...\n", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
