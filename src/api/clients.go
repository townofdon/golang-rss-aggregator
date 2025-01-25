package api

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/townofdon/tutorial-go-rss-server/internal/database"

	_ "github.com/lib/pq"
)

// singleton
var dbConn *sql.DB

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err.Error())
	}

	setupDBConnection()
}

func setupDBConnection() {
	if dbConn == nil {
		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL env var required")
		}

		var err error
		dbConn, err = sql.Open("postgres", dbUrl)

		if err != nil {
			log.Fatalf("Can't connect to database; using connection %v", dbUrl)
		}
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var required")
	}
	return port
}

func GetDBClient() (*database.Queries, *sql.DB) {
	setupDBConnection()

	db := database.New(dbConn)

	return db, dbConn
}

func WithDBClient(thunk func(db *database.Queries)) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env var required")
	}

	dbConn, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Can't connect to database")
	}

	defer dbConn.Close()

	thunk(database.New(dbConn))
}
