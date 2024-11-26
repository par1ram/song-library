package main

import (
	"time"
)

type InsertSongParams struct {
	GroupName string `json:"group_name" example:"Джизус"`
	SongName  string `json:"song_name" example:"Spirit of the world"`
}

type UpdateSongParams struct {
	GroupName   string `json:"group_name" example:"Джизус"`
	SongName    string `json:"song_name" example:"Spirit of the world"`
	Text        string `json:"text,omitempty" example:"Lyrics of the song"`
	ReleaseDate string `json:"release_date,omitempty" example:"2022-11-11"`
	Link        string `json:"link,omitempty" example:"https://example.com"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"An error occurred"`
}

type EmptyResponse struct{}

type ParametersGetSongs struct {
	GroupName string `json:"group_name" example:"Джизус"`
	SongName  string `json:"song" example:"Spirit of the world"`
	Limit     int32  `json:"limit" example:"10"`
	Offset    int32  `json:"offset" example:"0"`
}

type ParametersGetSongsWithFiltersAndPagination struct {
	Column1     string `json:"group_name" example:"Джизус"`
	Column2     string `json:"song" example:"Spirit of the world"`
	ReleaseDate string `json:"release_date" example:"2022-01-01"`
	Limit       int32  `json:"limit" example:"10"`
	Offset      int32  `json:"offset" example:"0"`
}

type Song struct {
	ID          string     `json:"id" example:"f7d5a812-12f3-4b92-a7d6-5c9f3a327f6b"`
	GroupName   string     `json:"group_name" example:"Джизус"`
	SongName    string     `json:"song_name" example:"Spirit of the world"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Text        string     `json:"text,omitempty" example:"Song lyrics"`
	Link        string     `json:"link,omitempty" example:"https://example.com"`
}

type Verse struct {
	ID     string `json:"id" example:"f7d5a812-12f3-4b92-a7d6-5c9f3a327f6b"`
	Limit  string `json:"limit" example:"10"`
	Offset string `json:"offset" example:"0"`
}
