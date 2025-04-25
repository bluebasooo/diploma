package entity

import "time"

type VideoShortPreview struct {
	ID          string `bson	:"_id"` // same as file meta
	Img         string `bson:"img"`
	Name        string `bson:"name"`
	DurationSec int64  `bson:"durationSec"`
	FileId      string `bson:"fileId"` // на всякий случай
}

type AuthorShortPreview struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Img  string `bson:"img"`
}

// TODO проработать - сейчас все что ниже накид на страницу самого видео
type VideoPreview struct {
	ID            string        `bson:"_id"`
	Name          string        `bson:"name"`
	Img           string        `bson:"img"`
	Description   string        `bson:"description"`
	AuthorPreview AuthorPreview `bson:"author"`
	VideoStats    VideoStats    `bson:"stats"`
	CreatedAt     time.Time     `bson:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt"`
	DurationMs    int64         `bson:"durationMs"`
	CommentTreeId string        `bson:"commentTreeId"`
}

type AuthorPreview struct {
	ID      string `bson:"_id"`
	Name    string `bson:"name"`
	Img     string `bson:"img"`
	SubsNum int64  `bson:"subsNum"`
}

type VideoStats struct {
	Views    int64 `bson:"views"`
	Likes    int64 `bson:"likes"`
	Dislikes int64 `bson:"dislikes"`
}

type Comment struct{}
