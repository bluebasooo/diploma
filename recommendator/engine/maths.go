package engine

import (
	"maps"
	"math"
)

func dist(vector map[string]float64) float64 {
	var dist float64
	for _, v := range vector {
		dist += v * v
	}
	return math.Sqrt(dist)
}

func diff(vector map[string]float64, otherVector map[string]float64) map[string]float64 {
	diff := maps.Clone(vector)
	for k, v := range otherVector {
		_, ok := diff[k]
		if !ok {
			diff[k] = 0.0
		}
		diff[k] -= v
	}
	return diff
}
