package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sort"
)

func ShouldAddToBucket(bucket *entity.Bucket, dotToAdd *entity.DotHistory) (float64, bool) {
	if len(bucket.BucketDotsToDistToCenter) <= 1 {
		return 0, true
	}

	if bucket.IsSeparated {
		return 0, false
	}

	distance := distToCenter(bucket.BucketCenter, dotToAdd)

	anomalyDotId := oneAnomalyDot(bucket)

	val := bucket.BucketDotsToDistToCenter[anomalyDotId]

	return distance, distance <= val
}

func AddDot(bucketId string, dot *entity.DotHistory) (*entity.Bucket, error) {
	var bucket *entity.Bucket
	for {
		if BucketRepo.LockOnBucket(bucketId) {
			bucket, _ = BucketRepo.GetBucket(context.Background(), bucketId)
			break
		}
	}
	defer BucketRepo.UnlockOnBucket(bucketId)

	bucketDotsIds := Plain(
		bucket.BucketDotsToDistToCenter,
		func(dotId entity.VideoDotId, dist float64) entity.VideoDotId {
			return dotId
		})

	bucketDots, err := DotsRepo.GetDots(context.Background(), bucketDotsIds)
	if err != nil {
		return nil, err
	}

	// если точка уже существует просто обновляем инфу о ней в бакете
	exist := false
	for i, val := range bucketDots {
		if val.ID == dot.ID {
			bucketDots[i] = *dot
			exist = true
		}
	}
	if !exist {
		bucketDots = append(bucketDots, *dot)
	}

	center := calculateCenter(bucketDots)
	bucket.BucketCenter = center

	distsToCenter := recalculateDistsToCenter(center, bucketDots)

	renewedBucket := &entity.Bucket{
		ID:                       bucket.ID,
		BucketDotsToDistToCenter: distsToCenter,
		BucketCenter:             center,
	}

	err = BucketRepo.UpsertBuckets(context.Background(), *renewedBucket)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return renewedBucket, nil
}

func recalculateDistsToCenter(center map[string]float64, dots []entity.DotHistory) map[string]float64 {
	distsToCenter := make(map[string]float64)
	for _, dot := range dots {
		distsToCenter[dot.GetDotID()] = distToCenter(center, &dot)
	}
	return distsToCenter
}

func oneAnomalyDot(bucket *entity.Bucket) string {
	var anomalyDotId string
	maxDist := 0.0
	for dotId, dist := range bucket.BucketDotsToDistToCenter {
		if dist > maxDist {
			maxDist = dist
			anomalyDotId = dotId
		}
	}
	return anomalyDotId
}

// можем разбить бакет на 2 когда, 2 корректных выброса имеют diff больше maxAngle
// первый возвращаемое значение - точка по которой делим
// второй - true если можно делить, false если нельзя
// третий - все точки отсортированные по расстоянию до точки по которой делим
func canSplitBucketOver(dotsFromBucket []entity.DotHistory, firstDot *entity.DotHistory) (bool, []Pair[string, float64], []Pair[string, float64]) {
	firstsDists := make([]Pair[string, float64], 0, len(dotsFromBucket))

	for _, dot := range dotsFromBucket {
		if dot.GetDotID() == firstDot.GetDotID() {
			continue
		}
		distOverSplitted := distBetweenDots(firstDot, &dot)
		firstsDists = append(firstsDists, Pair[string, float64]{Key: dot.GetDotID(), Value: distOverSplitted})
	}

	sort.Slice(firstsDists, func(i, j int) bool {
		return firstsDists[i].Value < firstsDists[j].Value
	})

	index, maxGrows := maxGrow(firstsDists)

	return maxGrows >= 100, firstsDists[:index], firstsDists[index:]
}

func ProcessSplitBucket(bucket *entity.Bucket) (*entity.Bucket, *entity.Bucket) {
	anomalyId := oneAnomalyDot(bucket)
	anomalyDot, err := DotsRepo.GetDot(context.Background(), anomalyId)
	if err != nil {
		return nil, nil
	}
	dotIdsFromBucket := Plain(
		bucket.BucketDotsToDistToCenter,
		func(dotId entity.VideoDotId, dist float64) entity.VideoDotId {
			return dotId
		})

	dotsFromBucket, err := DotsRepo.GetDots(context.Background(), dotIdsFromBucket)
	if err != nil {
		return nil, nil
	}

	ok, firstDotPairs, secondDotPairs := canSplitBucketOver(dotsFromBucket, anomalyDot)
	if !ok {
		return nil, nil
	}

	firstDotsIds := Map(
		firstDotPairs,
		func(pair Pair[string, float64]) entity.VideoDotId {
			return pair.Key
		})

	secondDotsIds := Map(
		secondDotPairs,
		func(pair Pair[string, float64]) entity.VideoDotId {
			return pair.Key
		})

	firstDots, err := DotsRepo.GetDots(context.Background(), firstDotsIds)
	if err != nil {
		return nil, nil
	}
	secondDots, err := DotsRepo.GetDots(context.Background(), secondDotsIds)
	if err != nil {
		return nil, nil
	}

	firstBucket := proceessCreateBucket(firstDots)
	secondBucket := proceessCreateBucket(secondDots)

	return firstBucket, secondBucket
}

func proceessCreateBucket(dots []entity.DotHistory) *entity.Bucket {
	center := calculateCenter(dots)
	distsToCenter := make(map[string]float64)
	for _, dot := range dots {
		distsToCenter[dot.GetDotID()] = distToCenter(center, &dot)
	}

	return &entity.Bucket{
		ID:                       primitive.NewObjectID().Hex(),
		BucketCenter:             center,
		BucketDotsToDistToCenter: distsToCenter,
	}
}
