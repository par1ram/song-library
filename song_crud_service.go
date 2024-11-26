package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/par1ram/song-library/internal/database"
)

// insertSong godoc
// @Summary Insert a new song
// @Description Inserts a new song into the database. Fetches additional song details from an external API.
// @Tags Songs
// @Accept json
// @Produce json
// @Param body body InsertSongParams true "Request body for inserting a song"
// @Success 201 {object} InsertSongParams "Created song details"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing fields"
// @Failure 500 {object} ErrorResponse "Failed to insert song or fetch external API data"
// @Router /songs/add [post]
func (apiCfg *apiConfig) insertSong(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupName string `json:"group_name"`
		SongName  string `json:"song_name"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Couldn't decode parameters", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.GroupName == "" || params.SongName == "" {
		RespondWithError(w, http.StatusBadRequest, "Group name and Song name cannot be empty")
		return
	}

	// --- Запрос к внешнему API ---
	type SongDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	apiURL := getExternalApiURL()
	query := url.Values{}
	query.Add("group", params.GroupName)
	query.Add("song", params.SongName)
	fullURL := apiURL + "?" + query.Encode()

	resp, err := apiCfg.HTTPClient.Get(fullURL)
	if err != nil {
		log.Println("Error making request to external API:", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to get song info")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("External API returned non-200 status:", resp.Status)
		RespondWithError(w, http.StatusInternalServerError, "Failed to get song info")
		return
	}

	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Println("Error decoding response:", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to process song info")
		return
	}
	// --- Конец запроса к внешнему API ---

	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		log.Println("Invalid date format from external API:", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to parse date")
		return
	}

	err = apiCfg.DB.InsertSong(r.Context(), database.InsertSongParams{
		ID:          uuid.New(),
		GroupName:   params.GroupName,
		SongName:    params.SongName,
		ReleaseDate: sql.NullTime{Time: releaseDate, Valid: true},
		Text:        sql.NullString{String: songDetail.Text, Valid: true},
		Link:        sql.NullString{String: songDetail.Link, Valid: true},
	})

	if err != nil {
		log.Println("Failed to insert song:", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to insert song")
		return
	}

	RespondWithJSON(w, 201, params)
}

// updateSong godoc
// @Summary Update an existing song
// @Description Updates the details of an existing song in the database.
// @Tags Songs
// @Accept json
// @Produce json
// @Param id query string true "Song ID (UUID)"
// @Param body body UpdateSongParams true "Request body for updating a song"
// @Success 200 {object} UpdateSongParams "Updated song details"
// @Failure 400 {object} ErrorResponse "Invalid request body or song ID"
// @Failure 500 {object} ErrorResponse "Failed to update song"
// @Router /songs/update [put]
func (apiCfg *apiConfig) updateSong(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupName   string `json:"group_name"`
		SongName    string `json:"song_name"`
		Text        string `json:"text"`
		ReleaseDate string `json:"release_date"`
		Link        string `json:"link"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Couldnt decode parameters", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.GroupName == "" || params.SongName == "" {
		RespondWithError(w, http.StatusBadRequest, "Group name and Song name cannot be empty")
		return
	}

	Text := sql.NullString{String: params.Text, Valid: params.Text != ""}
	Link := sql.NullString{String: params.Link, Valid: params.Link != ""}

	releaseDate, err := time.Parse("2006-01-02", params.ReleaseDate)
	if err != nil {
		log.Println("Date parsing error:", err)
		RespondWithError(w, http.StatusBadRequest, "Wrong date format. Use YYYY-MM-DD.")
		return
	}
	ReleaseDate := sql.NullTime{Time: releaseDate, Valid: params.ReleaseDate != ""}

	songIDStr := r.URL.Query().Get("id")
	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		log.Println(err)
		RespondWithError(w, 400, "Wrong song ID format.")
		return
	}

	err = apiCfg.DB.UpdateSong(r.Context(), database.UpdateSongParams{
		ID:          songID,
		GroupName:   params.GroupName,
		SongName:    params.SongName,
		Text:        Text,
		ReleaseDate: ReleaseDate,
		Link:        Link,
	})

	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to update song")
		return
	}

	RespondWithJSON(w, 200, params)
}

// deleteSong godoc
// @Summary Delete a song
// @Description Deletes a song from the database by its ID.
// @Tags Songs
// @Accept json
// @Produce json
// @Param id query string true "Song ID (UUID)"
// @Success 200 {object} EmptyResponse "Successful deletion"
// @Failure 400 {object} ErrorResponse "Invalid or missing song ID"
// @Failure 500 {object} ErrorResponse "Failed to delete song"
// @Router /songs/delete [delete]
func (apiCfg *apiConfig) deleteSong(w http.ResponseWriter, r *http.Request) {
	songIDStr := r.URL.Query().Get("id")
	if songIDStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing song ID")
		return
	}

	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		log.Println(err)
		RespondWithError(w, 400, "Invalid song ID")
		return
	}

	err = apiCfg.DB.DeleteSong(r.Context(), songID)
	if err != nil {
		log.Println(err)
		RespondWithError(w, 500, "Failed to delete song")
		return
	}

	RespondWithJSON(w, 200, struct{}{})
}
