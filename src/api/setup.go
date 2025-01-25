package api

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/townofdon/tutorial-go-rss-server/internal/database"

	_ "github.com/lib/pq"
)

func SetupApiClients() (clients Clients, port string) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err.Error())
	}

	port = os.Getenv("PORT")
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

	clients = Clients{
		DB: database.New(dbConn),
	}

	return
}
