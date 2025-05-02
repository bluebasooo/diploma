package utils

func ToMapByUniqueField[K comparable, V any](slice []V, keyFunc func(*V) K) map[K]V {
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		result[keyFunc(&item)] = item
	}
	return result
}
