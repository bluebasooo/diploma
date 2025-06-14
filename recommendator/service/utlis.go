package service

import (
	"cmp"
	"sort"
)

func GroupByValueProp[K comparable, V any](vals []V, keyExtractor func(V) K) map[K][]V {
	grouped := make(map[K][]V)
	for _, val := range vals {
		key := keyExtractor(val)
		if _, ok := grouped[key]; !ok {
			grouped[key] = make([]V, 0)
		}
		grouped[key] = append(grouped[key], val)
	}
	return grouped
}

func Plain[K comparable, V any, R any](vals map[K]V, extractor func(K, V) R) []R {
	res := make([]R, 0, len(vals))
	for k, v := range vals {
		res = append(res, extractor(k, v))
	}
	return res
}

func Map[V any, R any](vals []V, keyExtractor func(V) R) []R {
	res := make([]R, 0, len(vals))
	for _, val := range vals {
		key := keyExtractor(val)
		res = append(res, key)
	}
	return res
}

func sortedByValueKeys[K comparable, V comparable](m map[K]V, comparator func(one V, two V) bool) []K {
	pairs := toPairs(m)
	sort.Slice(pairs, func(i int, j int) bool {
		return comparator(pairs[i].Value, pairs[j].Value)
	})

	return keysFromPairs(pairs)
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

func toPairs[K comparable, V any](m map[K]V) []Pair[K, V] {
	result := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, Pair[K, V]{Key: k, Value: v})
	}
	return result
}

type Pair[K any, V any] struct {
	Key   K
	Value V
}

func keysFromPairs[K any, V any](pairs []Pair[K, V]) []K {
	result := make([]K, 0, len(pairs))
	for _, pair := range pairs {
		result = append(result, pair.Key)
	}
	return result
}
