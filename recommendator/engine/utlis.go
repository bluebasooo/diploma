package engine

import (
	"cmp"
	"math"
	"sort"
)

func values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func diffBetween(a []float64) []float64 {
	result := make([]float64, 0, len(a))
	for i := 0; i < len(a)-1; i++ {
		result = append(result, a[i+1]-a[i])
	}
	return result
}

func maxOver[K cmp.Ordered](a []K) K {
	max, _ := maxOverWithIndex(a)
	return max
}

func maxDiffOverPairsWithIndex[K cmp.Ordered](a []Pair[K, float64]) (Pair[K, float64], float64) {
	max := math.Inf(-1)
	index := -1

	for i := 1; i < len(a); i++ {
		if a[i].Value-a[i-1].Value > max {
			max = a[i].Value - a[i-1].Value
			index = i
		}
	}

	return a[index], max
}

func maxGrow[K cmp.Ordered](a []Pair[K, float64]) (int, float64) {
	if len(a) < 2 {
		return 0, 1.0
	}

	max := a[1].Value - a[0].Value
	prev := a[1].Value - a[0].Value
	grow := 1.0
	index := -1

	for i := 2; i < len(a); i++ {
		prev = a[i].Value - a[i-1].Value
		if prev > max {
			max = prev
			index = i
			grow = max / prev
		}
	}

	return index, grow
}

func maxOverWithIndex[K cmp.Ordered](a []K) (K, int) {
	max := a[0]
	index := 0
	for i, v := range a {
		if v > max {
			max = v
			index = i
		}
	}
	return max, index
}

func Pairs[K any, V any](m map[K]V) []Pair[K, V] {
	result := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Pair[K, V]{Key: k, Value: v})
	}
	return result
}

func SortedPairsByValue[K comparable, V comparable](m map[K]V, less func(a, b V) bool) []Pair[K, V] {
	pairs := Pairs(m)

	sort.Slice(pairs, func(i, j int) bool {
		return less(pairs[i].Value, pairs[j].Value)
	})
	return pairs
}

func defaultLess[K cmp.Ordered](a, b K) bool {
	return a < b
}

func reverseLess[K cmp.Ordered](a, b K) bool {
	return a > b
}

type Pair[K any, V any] struct {
	Key   K
	Value V
}

func RangedKeysPairs[K any, V any](pairs []Pair[K, V], start int, end int) []K {
	if start < 0 {
		start = 0
	}

	if end > len(pairs) {
		end = len(pairs)
	}

	return KeysFromPairs(pairs[start:end])
}

func KeysFromPairs[K any, V any](pairs []Pair[K, V]) []K {
	result := make([]K, 0, len(pairs))
	for _, pair := range pairs {
		result = append(result, pair.Key)
	}
	return result
}

func toPairs[K comparable, V any](m map[K]V) []Pair[K, V] {
	result := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Pair[K, V]{Key: k, Value: v})
	}
	return result
}

func toMap[K comparable, V any](pairs []Pair[K, V]) map[K]V {
	result := make(map[K]V, len(pairs))
	for _, pair := range pairs {
		result[pair.Key] = pair.Value
	}
	return result
}
