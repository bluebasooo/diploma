package dto

import "time"

type HistoryDto struct {
	UserID    string
	VideoID   string
	CreatedAt time.Time
	Metric    float64
}
