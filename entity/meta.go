package entity

import "time"

type FileMeta struct {
	ID        string    `json:"id"`
	PartsIDs  []string  `json:"parts"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VideoShortPreview struct {
	ID          string `json:"id"`
	Img         string `json:"img"`
	Name        string `json:"name"`
	DurationSec int64  `json:"durationSec"`
}

type AuthorShortPreview struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

// TODO проработать - сейчас все что ниже накид на страницу самого видео
type VideoPreview struct {
	ID            string     `json:"id"`
	Description   string     `json:"description"`
	AuthorPreview string     `json:"author"`
	VideoStats    VideoStats `json:"stats"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type VideoStats struct {
	Views       int64 `json:"views"`
	Likes       int64
	Dislikes    int64
	CommentsNum int64
}

type Comment struct {
}
