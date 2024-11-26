package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/par1ram/song-library/internal/database"
)

type apiConfig struct {
	DB         *database.Queries
	HTTPClient *http.Client
}

func createApiConfig(con *sql.DB) apiConfig {
	apiCfg := apiConfig{
		DB:         database.New(con),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	return apiCfg
}

func connectToDatabase() *sql.DB {
	DB_URL := getDatabaseURL()

	connection, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	return connection
}
