package entity

import "time"

// only metrics from viewing
type Metric struct {
	ID        string    `sql:"id"`
	UserID    string    `sql:"user_id"`
	VideoID   string    `sql:"video_id"`
	Type      string    `sql:"type"`
	Value     float64   `sql:"value"`
	CreatedAt time.Time `sql:"created_at"`
}
