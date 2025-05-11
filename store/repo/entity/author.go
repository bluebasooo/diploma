package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Author struct {
	ID                  *primitive.ObjectID `bson:"_id"`
	Username            string              `bson:"username"`
	AuthorStats         AuthorStats         `bson:"author_stats"`
	VideoIDs            []string            `bson:"video_ids"`
	AuthorSubscriptions []string            `bson:"author_subscriptions"`
	ImgLink             string              `bson:"img_link"`
}

type AuthorStats struct {
	Subscribers int64 `bson:"subscribers"`
	Likes       int64 `bson:"likes"`
	Views       int64 `bson:"views"`
}
