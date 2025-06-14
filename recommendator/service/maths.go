package service

import (
	"dev/bluebasooo/video-recomendator/entity"
	"math"
)

func distToCenter(center map[string]float64, dot *entity.DotHistory) float64 {
	diff := make(map[string]float64)

	for k, v := range center {
		diff[k] = dot.GetValue(k) - v
	}

	dist := 0.0
	for _, v := range diff {
		dist += v * v
	}
	return math.Sqrt(dist)
}

func distBetweenDots(dot *entity.DotHistory, otherDot *entity.DotHistory) float64 {
	diff := make(map[string]float64)

	for k, v := range otherDot.History {
		if _, ok := diff[k]; !ok {
			diff[k] = 0.0
		}
		diff[k] -= v
	}

	for k, v := range dot.History {
		if _, ok := diff[k]; !ok {
			diff[k] = 0.0
		}
		diff[k] += v
	}

	dist := 0.0
	for _, v := range diff {
		dist += v * v
	}

	return math.Sqrt(dist)
}

// not empty vector list
func calculateCenter(dots []entity.DotHistory) map[string]float64 {
	if len(dots) == 0 {
		return make(map[string]float64)
	}

	sum := make(map[string]float64)
	for _, dot := range dots {
		for k, _ := range dot.History {
			if _, ok := sum[k]; !ok {
				sum[k] = 0.0
			}
			sum[k] += dot.GetValue(k)
		}
	}

	len := float64(len(dots))
	for k, v := range sum {
		sum[k] = v / len
	}
	return sum
}
