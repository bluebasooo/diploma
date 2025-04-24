package engine

import "cmp"

func values[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
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
	max := a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}
