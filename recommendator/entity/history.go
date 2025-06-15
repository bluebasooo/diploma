package entity

import "time"

// complete history value on viewing
type History struct {
	ID        string    `ch:"id"`
	UserID    string    `ch:"user_id"`
	VideoID   string    `ch:"video_id"`
	CreatedAt time.Time `ch:"created_at"`
	Metric    float64   `ch:"value"`
}

type UserHistory struct {
	UserID  string  `ch:"user_id"`
	VideoID string  `ch:"video_id"`
	Metric  float64 `ch:"value"`
}

// в pool хотелось бы ранжирование по metric + недавности просмотра

type ShortUserHistory struct {
	UserID    string    `ch:"user_id"`
	VideoID   string    `ch:"video_id"`
	CreatedAt time.Time `ch:"created_at"`
}
