package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/townofdon/tutorial-go-rss-server/internal/database"
	"github.com/townofdon/tutorial-go-rss-server/src/api"
	v1 "github.com/townofdon/tutorial-go-rss-server/src/api/v1"
	"github.com/townofdon/tutorial-go-rss-server/src/auth"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var required")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env var required")
	}

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	clients := api.Clients{
		DB: database.New(dbConn),
	}

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
	v1Router.Get("/users/current", auth.Handler(&clients, v1.GetUserByApiKey))
	v1Router.Post("/users", api.Handler(&clients, v1.CreateUser))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on port %s...\n", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
