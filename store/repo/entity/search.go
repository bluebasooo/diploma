package entity

type VideoIndex struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DurationMs  int    `json:"durationMs"`
	AuthorName  string `json:"authorName"`
	UploadDate  string `json:"uploadDate"`
	Views       int    `json:"views"`
	Hidden      bool   `json:"hidden"`
}

type AuthorIndex struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type VideoSearchResult struct {
	VideoIndex
	Searchable
}

type AuthorSearchResult struct {
	AuthorIndex
	Searchable
}

type Searchable struct {
	ID    string  `json:"_id"`
	Score float64 `json:"_score"`
	Index string  `json:"_index"`
}
