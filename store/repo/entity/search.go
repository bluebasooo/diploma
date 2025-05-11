package entity

import "time"

type VideoIndex struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DurationMs  int       `json:"durationMs"`
	AuthorName  string    `json:"authorName"`
	UploadDate  time.Time `json:"uploadDate"`
	Views       int       `json:"views"`
	Hidden      bool      `json:"hidden"`
}

type AuthorIndex struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type VideoSearchResult struct {
	Searchable
	VideoIndex
}

type AuthorSearchResult struct {
	Searchable
	AuthorIndex
}

type Searchable struct {
	ID    string  `json:"_id"`
	Score float64 `json:"_score"`
}
