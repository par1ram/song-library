package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getPort() string {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Couldnt get port from env")
	}

	return PORT
}

func getDatabaseURL() string {
	godotenv.Load()

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("Couldnt get port from env")
	}

	return DB_URL
}

func getExternalApiURL() string {
	godotenv.Load()

	EXTERNAL_API_URL := os.Getenv("EXTERNAL_API_URL")
	if EXTERNAL_API_URL == "" {
		log.Fatal("Couldnt get port from env")
	}

	return EXTERNAL_API_URL
}
