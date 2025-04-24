package entity

import "time"

type VideoShortPreview struct {
	ID          string `json:"id"` // same as file meta
	Img         string `json:"img"`
	Name        string `json:"name"`
	DurationSec int64  `json:"durationSec"`
	FileId      string `json:"fileId"` // на всякий случай
}

type AuthorShortPreview struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

// TODO проработать - сейчас все что ниже накид на страницу самого видео
type VideoPreview struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Img           string        `json:"img"`
	Description   string        `json:"description"`
	AuthorPreview AuthorPreview `json:"author"`
	VideoStats    VideoStats    `json:"stats"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	DurationMs    int64         `json:"durationMs"`
	CommentTreeId string        `json:"commentTreeId"`
}

type AuthorPreview struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Img     string `json:"img"`
	SubsNum int64  `json:"subsNum"`
}

type VideoStats struct {
	Views    int64 `json:"views"`
	Likes    int64
	Dislikes int64
}

type Comment struct{}
