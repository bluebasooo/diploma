package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToMapByUniqueField[K comparable, V any](slice []V, keyFunc func(*V) K) map[K]V {
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		result[keyFunc(&item)] = item
	}
	return result
}

func GetRandomHashId() string {
	return primitive.NewObjectID().Hex()
}

func NormalizeIds(ids []string) []string {
	setIds := make(map[string]bool, len(ids))

	for _, id := range ids {
		setIds[id] = true
	}

	normalizedAuthorIds := make([]string, 0, len(setIds))
	for k, _ := range setIds {
		normalizedAuthorIds = append(normalizedAuthorIds, k)
	}

	return normalizedAuthorIds
}
