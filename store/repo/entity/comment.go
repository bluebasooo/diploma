package entity

import "time"

type Comment struct {
	ID        string        `bson:"_id"`
	VideoID   string        `bson:"videoId"`
	Message   string        `bson:"text"` // TODO make limit on message
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
	Likes     int           `bson:"likes"`
	Dislikes  int           `bson:"dislikes"`
	Author    CommentAuthor `bson:"author"`
	ParentID  string        `bson:"parentId"`
	RootID    string        `bson:"rootId"`
	Relevance float64       `bson:"relevance"`
}

type CommentAuthor struct {
	ID       string `bson:"author_id"`
	Username string `bson:"username"`
}
