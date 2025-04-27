package engine

import "time"

type DotVideoMetrics struct {
	UserID     string
	Views      int
	Likes      int
	Dislikes   int
	WatchTime  float64 // in minutes
	ShareCount int
	CreatedAt  time.Time
	VideoID    string
	ViewID     string
}

func calculateMetric(videoMetrics DotVideoMetrics) float64 {
	return float64(videoMetrics.Views)*0.1 + float64(videoMetrics.Likes)*0.2 + float64(videoMetrics.Dislikes)*0.2 + videoMetrics.WatchTime*0.3 + float64(videoMetrics.ShareCount)*0.4
}
