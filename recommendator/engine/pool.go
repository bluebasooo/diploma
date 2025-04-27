package engine

type Pool struct {
	ID             string
	BucketID       string
	RangedVideoIDs []Pair[string, float64] // sorted by metrics
}

// mocked pool
func GetPool(id string) *Pool {
	pool := &Pool{
		ID: id,
	}

	return pool
}

func recalculatePool(pool *Pool) *Pool {
	bucket := getBucket(pool.BucketID)

	dotIds := keys(bucket.DotsToDistToCenter)

	history := getHistories(dotIds)

	videoIds := normalizeVideoIdsFromHistory(history)

	pool.RangedVideoIDs = videoIds

	return pool
}

func normalizeVideoIdsFromHistory(history []History) []Pair[string, float64] {
	videoIds := make(map[string]float64)
	for _, h := range history {
		for videoId, metrics := range h.Videos {
			if _, ok := videoIds[videoId]; !ok {
				videoIds[videoId] = 0
			}
			videoIds[videoId] += metrics
		}
	}

	return SortedPairsByValue(videoIds, reverseLess[float64])
}

// ранжирование путем ближайших пользователей + их видеопросмотре
