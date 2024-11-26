package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/par1ram/song-library/docs"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Song Library API
// @version 1.0
// @description API for managing a song library.

// @host localhost:8000
// @BasePath /

func main() {
	PORT := getPort()

	db := connectToDatabase()
	apiCfg := createApiConfig(db)

	if err := goose.Up(db, "sql/schema"); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/swagger/doc.json")))

	router.Post("/songs/getAll-with-filters", apiCfg.getSongs)
	router.Post("/songs/get-with-filters-and-pagination", apiCfg.getSongWithFiltersAndPagination)
	router.Get("/songs/get-song-verses", apiCfg.getSongVersesWithPagination)

	router.Post("/songs/add", apiCfg.insertSong)
	router.Put("/songs/update", apiCfg.updateSong)
	router.Delete("/songs/delete", apiCfg.deleteSong)

	server := &http.Server{
		Addr:           ":" + PORT,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1mb
	}

	fmt.Println("Server started on port:", PORT)
	log.Fatal(server.ListenAndServe())
}
