package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type MetricDto struct {
	UserID    string    `json:"userID"`
	VideoID   string    `json:"videoID"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	ViewID    string    `json:"viewID"`
}

var (
	videoPool = []string{
		"typeA_0", "typeA_1", "typeA_2", "typeA_3", "typeA_4",
		"typeA_5", "typeA_6", "typeA_7", "typeA_8", "typeA_9",
		"typeB_0", "typeB_1", "typeB_2", "typeB_3", "typeB_4",
		"typeB_5", "typeB_6", "typeB_7", "typeB_8", "typeB_9",
	}
	metricTypes = []string{"LIKE", "DISLIKE", "WATCH_TIME", "SHARE"}
)

func Generate() {
	rand.Seed(time.Now().UnixNano())
	var metrics []MetricDto
	now := time.Now()

	for userIdx := 0; userIdx < 100; userIdx++ {
		userID := fmt.Sprintf("user_%d", userIdx)
		viewCount := rand.Intn(6) + 5 // 5–10 views

		for viewIdx := 0; viewIdx < viewCount; viewIdx++ {
			viewID := fmt.Sprintf("view_%d_%d", userIdx, viewIdx)
			videoID := videoPool[rand.Intn(len(videoPool))]

			createdAt := now.Add(time.Duration(userIdx*100+viewIdx) * time.Second)

			// START
			metrics = append(metrics, MetricDto{
				UserID:    userID,
				VideoID:   videoID,
				Type:      "START",
				Value:     0,
				CreatedAt: createdAt,
				ViewID:    viewID,
			})

			// Добавляем рандомные внутренние метрики (только если value != 0)
			numInnerMetrics := rand.Intn(4) + 1 // от 1 до 4
			for i := 0; i < numInnerMetrics; i++ {
				mType := metricTypes[rand.Intn(len(metricTypes))]
				value := generateValue(mType)
				if value == 0 {
					continue // не включаем метрику с value == 0
				}
				metrics = append(metrics, MetricDto{
					UserID:    userID,
					VideoID:   videoID,
					Type:      mType,
					Value:     value,
					CreatedAt: createdAt.Add(time.Duration(i+1) * time.Second),
					ViewID:    viewID,
				})
			}
			last := metrics[len(metrics)-1]
			metrics = append(metrics, MetricDto{
				UserID:    userID,
				VideoID:   videoID,
				Type:      "END",
				Value:     0,
				CreatedAt: last.CreatedAt.Add(time.Duration(1) * time.Second),
				ViewID:    viewID,
			})

		}
	}

	// Сохраняем в JSON
	file, err := os.Create("metrics.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metrics); err != nil {
		panic(err)
	}

	fmt.Println("✅ Метрики успешно записаны в metrics.json")
}

func generateValue(metricType string) float64 {
	switch metricType {
	case "LIKE":
		return float64(rand.Intn(2)) // 0 или 1
	case "DISLIKE":
		return float64([]int{-1, 0}[rand.Intn(2)])
	case "WATCH_TIME":
		return float64(rand.Intn(1000) + 1) // 1–1000
	case "SHARE":
		return float64(rand.Intn(10) + 1) // 1–10
	default:
		return 0
	}
}
