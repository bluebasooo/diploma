package dto

import "time"

type MetricDto struct {
	UserID    string    `json:"userID"`
	VideoID   string    `json:"videoID"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	ViewID    string    `json:"viewID"`
}

// maybe grpc
type MetricsDto struct {
	Metrics []MetricDto `json:"metrics"`
}
