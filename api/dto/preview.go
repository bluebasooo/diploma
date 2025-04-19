package dto

type VideoPreviewDto struct { // on page of video
	Name       string                `json:"name"`
	Author     AuthorShortPreviewDto `json:"author"`
	DurationMs int64                 `json:"durationMs"`
	Stats      VideoStatsDto         `json:"stats"`
}

type VideoShortPreview struct { // on main page
	ID          string `json:"id"`
	Img         string `json:"img"`
	Name        string `json:"name"`
	DurationSec int64  `json:"durationSec"`
}

type AuthorShortPreviewDto struct { // all of them
	ID   string `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

type VideoStatsDto struct {
	Views      int64 `json:"views"`
	Likes      int64 `json:"likes"`
	Dislikes   int64 `json:"dislikes"`
	CommentsId int64 `json:"commentsId"`
}
