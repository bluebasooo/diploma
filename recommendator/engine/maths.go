package engine

import (
	"maps"
	"math"
)

func distBetweenDots(dot *Dot, otherDot *Dot) float64 {
	diff := diff(dot.History, otherDot.History)

	return dist(diff)
}

func dist(vector map[string]float64) float64 {
	var dist float64
	for _, v := range vector {
		dist += v * v
	}
	return math.Sqrt(dist)
}

func scalarProduct(vector map[string]float64, otherVector map[string]float64) float64 {
	scalar := 0.0
	for k, v := range vector {
		if _, ok := otherVector[k]; !ok {
			continue
		}
		scalar += v * otherVector[k]
	}
	return scalar
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

func plus(vector map[string]float64, otherVector map[string]float64) map[string]float64 {
	plus := maps.Clone(vector)
	for k, v := range otherVector {
		plus[k] += v
	}
	return plus
}

func dirtyMultiply(vector map[string]float64, num float64) map[string]float64 {
	for k, v := range vector {
		vector[k] = v * num
	}
	return vector
}

func divide(vector map[string]float64, num float64) map[string]float64 {
	divided := maps.Clone(vector)
	for k, v := range divided {
		divided[k] = v / num
	}
	return divided
}

func dirtyDivide(vector map[string]float64, num float64) map[string]float64 {
	for k, v := range vector {
		vector[k] = v / num
	}
	return vector
}
