package main

import (
	"bytes"
	"dev/bluebasooo/rec-data-generator/generator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var needGenerate = false

func main() {
	if needGenerate {
		generator.Generate()
	}

	send()
}

type MetricsDto struct {
	Metrics []generator.MetricDto `json:"metrics"`
}

func send() {
	fileBytes, err := ioutil.ReadFile("metrics.json")
	if err != nil {
		panic(err)
	}

	var allMetrics []generator.MetricDto
	if err := json.Unmarshal(fileBytes, &allMetrics); err != nil {
		panic(err)
	}

	batchSize := 40
	url := "http://localhost:8080/metrics/"

	for i := 0; i < len(allMetrics); i += batchSize {
		end := i + batchSize
		if end > len(allMetrics) {
			end = len(allMetrics)
		}

		batch := allMetrics[i:end]
		payload := MetricsDto{Metrics: batch}

		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Ошибка маршалинга:", err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Ошибка запроса:", err)
			continue
		}

		resp.Body.Close()

		fmt.Printf("Отправлено %d метрик, статус: %s\n", len(batch), resp.Status)
	}
}
