package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/par1ram/song-library/internal/database"
)

// getSongs godoc
// @Summary Get all songs with filters
// @Description Get a list of songs with filtering and pagination options
// @Tags Songs
// @Accept json
// @Produce json
// @Param body body ParametersGetSongs true "Request body for getting songs with filters"
// @Success 200 {array} Song "List of songs"
// @Failure 400 {object} ErrorResponse "Invalid request body or filters"
// @Failure 500 {object} ErrorResponse "Failed to fetch songs"
// @Router /songs/getAll-with-filters [post]
func (apiCfg *apiConfig) getSongs(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupName string `json:"group_name"`
		SongName  string `json:"song"`
		Limit     int32  `json:"limit"`
		Offset    int32  `json:"offset"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Couldt decode parameters:", err)
		RespondWithError(w, 400, "Invalid request body")
		return
	}

	if params.GroupName == "" && params.SongName == "" {
		RespondWithJSON(w, 200, []database.Song{})
		return
	}

	if params.Limit < 0 || params.Offset < 0 {
		RespondWithError(w, 400, "Limit and Offset cannot be negative")
		return
	}
	const maxLimit = 100
	if params.Limit > maxLimit {
		params.Limit = maxLimit
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	songs, err := apiCfg.DB.GetSongs(r.Context(), database.GetSongsParams{
		GroupName: params.GroupName,
		SongName:  params.SongName,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, 404, "Song not found.")
			return
		}
		log.Println(err)
		RespondWithError(w, 500, "Couldnt get songs")
		return
	}

	RespondWithJSON(w, 200, songs)
}

// getSongWithFiltersAndPagination godoc
// @Summary Get song with filters and pagination
// @Description Get song by applying filters and pagination
// @Tags Songs
// @Accept json
// @Produce json
// @Param body body ParametersGetSongsWithFiltersAndPagination true "Request body for getting songs with filters and pagination"
// @Success 200 {array} Song "List of filtered songs"
// @Failure 400 {object} ErrorResponse "Invalid request body or filters"
// @Failure 500 {object} ErrorResponse "Failed to fetch songs"
// @Router /songs/get-with-filters-and-pagination [post]
func (apiCfg *apiConfig) getSongWithFiltersAndPagination(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Column1     string `json:"group_name"`
		Column2     string `json:"song"`
		ReleaseDate string `json:"release_date"`
		Limit       int32  `json:"limit"`
		Offset      int32  `json:"offset"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Couldnt decode parameters:", err)
		RespondWithError(w, 400, "Invalid request body")
	}

	Column1 := sql.NullString{String: params.Column1, Valid: params.Column1 != ""}
	Column2 := sql.NullString{String: params.Column2, Valid: params.Column2 != ""}

	var ReleaseDate sql.NullTime
	if params.ReleaseDate != "" {
		releaseDate, err := time.Parse("2006-01-02", params.ReleaseDate)
		if err != nil {
			RespondWithError(w, 400, "Wrong date format. Use YYYY-MM-DD.")
			log.Println("Date parsing error:", err)
			return
		}
		ReleaseDate = sql.NullTime{Time: releaseDate, Valid: true}
	}

	if params.Column1 == "" && params.Column2 == "" && params.ReleaseDate == "" {
		RespondWithError(w, http.StatusBadRequest, "At least one filter is required")
		return
	}

	if params.Limit < 0 || params.Offset < 0 {
		RespondWithError(w, http.StatusBadRequest, "Limit and Offset cannot be negative")
		return
	}
	const maxLimit = 100
	if params.Limit > maxLimit {
		params.Limit = maxLimit
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	song, err := apiCfg.DB.GetSongWithFiltersAndPagination(r.Context(), database.GetSongWithFiltersAndPaginationParams{
		Column1:     Column1,
		Column2:     Column2,
		ReleaseDate: ReleaseDate,
		Limit:       params.Limit,
		Offset:      params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, 404, "Song not found.")
			return
		}
		log.Println("Couldnt get songs:", err)
		RespondWithError(w, 500, "Couldnt get songs")
		return
	}

	RespondWithJSON(w, 200, song)
}

// getSongVersesWithPagination godoc
// @Summary Get song verses with pagination
// @Description Get song verses by song ID with pagination
// @Tags Songs
// @Accept json
// @Produce json
// @Param id query string true "Song ID (UUID)"
// @Param limit query int32 false "Limit for pagination"
// @Param offset query int32 false "Offset for pagination"
// @Success 200 {array} Verse "List of song verses"
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 404 {object} ErrorResponse "Song not found"
// @Failure 500 {object} ErrorResponse "Failed to get song verses"
// @Router /songs/get-song-verses [get]
func (apiCfg *apiConfig) getSongVersesWithPagination(w http.ResponseWriter, r *http.Request) {
	songIDStr := r.URL.Query().Get("id")
	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusBadRequest, "Wrong song ID format.")
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		log.Println(err)
		RespondWithError(w, 400, err.Error())
		return
	}

	verses, err := apiCfg.DB.GetSongVersesWithPagination(r.Context(), database.GetSongVersesWithPaginationParams{
		ID:     songID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			RespondWithError(w, 404, "Song not found.")
			return
		}
		log.Println(err)
		RespondWithError(w, 500, "Error getting verses.")
		return
	}

	RespondWithJSON(w, 200, verses)
}

// Вспомогательная функция для парсинга limit и offset из запроса
func parseLimitOffset(r *http.Request) (int32, int32, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10)
	if limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(parsedLimit)
		} else {
			return 0, 0, fmt.Errorf("wrong limit format: %s", limitStr)
		}
	}

	offset := int32(0)
	if offsetStr != "" {
		if parsedOffset, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(parsedOffset)
		} else {
			return 0, 0, fmt.Errorf("wrong offset format: %s", offsetStr)
		}
	}

	return limit, offset, nil
}
