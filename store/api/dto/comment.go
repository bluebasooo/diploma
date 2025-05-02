package dto

import "time"

type CommentDto struct {
	ID        string           `json:"id"`
	VideoID   string           `json:"videoId"`
	Message   string           `json:"message"`
	Likes     int              `json:"likes"`
	Dislikes  int              `json:"dislikes"`
	Author    CommentAuthorDto `json:"author"`
	Childs    []CommentDto     `json:"childs"`
	ParentID  string           `json:"parentId"`
	RootID    string           `json:"rootId"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type CommentAuthorDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
