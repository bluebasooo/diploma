package engine

import (
	"math"
	"sort"
)

const maxAngle = 0.1

func (dot *Dot) dotDist() float64 {
	return math.Sqrt(dist(dot.History))
}

func (dot *Dot) distToCenter(center map[string]float64) float64 {
	difference := diff(dot.History, center)
	return dist(difference)
}

func (dot *Dot) distFromDot(otherDot *Dot) float64 {
	difference := diff(dot.History, otherDot.History)
	return dist(difference)
}

func calculateCenter(dots []Dot) map[string]float64 {
	sum := make(map[string]float64)
	for _, dot := range dots {
		for k, v := range dot.History {
			_, ok := sum[k]
			if !ok {
				sum[k] = 0.0
			}
			sum[k] += v
		}
	}

	len := float64(len(dots))
	for k, v := range sum {
		sum[k] = v / len
	}
	return sum
}

func (bucket *Bucket) shouldAddToBucket(dotToAdd *Dot) bool {
	distance := dotToAdd.distToCenter(bucket.Center)

	dists := values(bucket.DotsToDistToCenter)

	sort.Float64Slice(dists).Sort()

	diffs := diffBetween(dists)
	maxOverDiffs := maxOver(diffs)

	return math.Abs(distance-maxOverDiffs) > maxAngle
}

func (bucket *Bucket) canDifferentiateBucket() bool {
	dists := values(bucket.DotsToDistToCenter)

	diffs := diffBetween(dists)
	maxOverDiffs := maxOver(diffs)

	return math.Abs(distance-maxOverDiffs) > maxAngle
}

// должна быть операция выделения в новый бакет
// нужно брать выбросы и на основе их выделять новый бакет
// можно делать на основе предпосчитанного центра
// должен пересчитываться на добавление точик

// должна быть операция перемещения в новый бакет
// если становится слишком большим выбросом
