package entity

type Author struct {
	ID                        int64       `bson:"_id"`
	Username                  string      `bson:"username"`
	AuthorStats               AuthorStats `bson:"author_stats"`
	VideoIDs                  []string    `bson:"video_ids"`
	SubscriptionsOnAuthorsIDs []string    `bson:"subscriptions_on_authors_ids"`
	ImgLink                   string      `bson:"img_link"`
}

type AuthorStats struct {
	Subscribers int64 `bson:"subscribers"`
	Likes       int64 `bson:"likes"`
	Views       int64 `bson:"views"`
}
