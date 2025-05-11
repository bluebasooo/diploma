package dto

import "time"

type VideoPreviewDto struct { // on page of video
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Author      AuthorShortPreviewDto `json:"author"`
	DurationMs  int64                 `json:"durationMs"`
	Stats       VideoStatsDto         `json:"stats"`
	CreatedAt   time.Time             `json:"createdAt"`
}

type VideoShortPreviewDto struct { // on main page
	ID         string                `json:"id"`
	Img        string                `json:"img"`
	Name       string                `json:"name"`
	DurationMs int64                 `json:"durationMs"`
	Author     AuthorShortPreviewDto `json:"author"`
	Views      int64                 `json:"views"`
}

type AuthorShortPreviewDto struct { // all of them
	Username string `json:"username"`
	Img      string `json:"img"`
	SubsNum  int64  `json:"subsNum"`
}

type VideoStatsDto struct {
	Views    int64 `json:"views"`
	Likes    int64 `json:"likes"`
	Dislikes int64 `json:"dislikes"`
}

type CreateVideoPreviewDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DurationMs  int64  `json:"durationMs"`
}
