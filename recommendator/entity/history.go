package entity

import "time"

// complete history value on viewing
type History struct {
	ID        string    `sql:"id"`
	UserID    string    `sql:"user_id"`
	VideoID   string    `sql:"video_id"`
	Metric    float64   `sql:"value"`
	CreatedAt time.Time `sql:"created_at"`
}
