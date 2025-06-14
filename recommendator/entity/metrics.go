package entity

import "time"

type MetricType string

const (
	MetricTypeUnknown   MetricType = "START"
	MetricTypeLike      MetricType = "LIKE"
	MetricTypeDislike   MetricType = "DISLIKE"
	MetricTypeShare     MetricType = "WATCH_TIME"
	MetricTypeWatchTime MetricType = "SHARE"
	MetricTypeStart     MetricType = "END"
	MetricTypeProcessed MetricType = "PROCESSED"
)

// only metrics from viewing
type Metric struct {
	ID        string    `sql:"id"`
	UserID    string    `sql:"user_id"`
	VideoID   string    `sql:"video_id"`
	ViewID    string    `sql:"view_id"`
	Type      string    `sql:"type"`
	Value     float64   `sql:"value"`
	CreatedAt time.Time `sql:"created_at"`
}

type ViewIdentifier struct {
	UserID  string
	VideoID string
	ViewID  string
}
