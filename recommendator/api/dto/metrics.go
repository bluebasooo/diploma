package dto

import "time"

type MetricDto struct {
	UserID    string
	VideoID   string
	Type      string
	Value     float64
	CreatedAt time.Time
	ViewID    string
}

// maybe grpc
type MetricsDto struct {
	Metrics []MetricDto `json:"metrics"`
}
