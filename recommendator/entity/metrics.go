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
	ID        string    `ch:"id"`
	UserID    string    `ch:"user_id"`
	VideoID   string    `ch:"video_id"`
	ViewID    string    `ch:"view_id"`
	Type      string    `ch:"type"`
	Value     float64   `ch:"value"`
	CreatedAt time.Time `ch:"created_at"`
}

type ViewIdentifier struct {
	ViewID  string `ch:"view_id"`
	UserID  string `ch:"user_id"`
	VideoID string `ch:"video_id"`
}
