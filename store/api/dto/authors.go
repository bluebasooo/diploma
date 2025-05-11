package dto

type CreateAuthorDto struct {
	Username string `json:"username"`
}

type AuthorDto struct {
	ID            string            `json:"id"`
	Username      string            `json:"username"`
	ImgLink       string            `json:"img_link"`
	Stats         StatsDto          `json:"stats"`
	Subscriptions []SubscriptionDto `json:"subscriptions"`
}

type StatsDto struct {
	Subscribers int64 `json:"subscribers"`
	Likes       int64 `json:"likes"`
	Views       int64 `json:"views"`
}

type SubscriptionDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	ImgLink  string `json:"img_link"`
}
