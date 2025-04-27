package engine

import (
	"sort"
)

type Bucket struct {
	ID                 string
	DotsToDistToCenter map[string]float64
	Center             map[string]float64
}

func getBucket(id string) *Bucket {
	return &Bucket{
		ID:                 id,
		DotsToDistToCenter: make(map[string]float64),
		Center:             make(map[string]float64),
	}
}

func (bucket *Bucket) shouldAddToBucket(dotToAdd *Dot) bool {
	if len(bucket.DotsToDistToCenter) == 1 {
		return true
	}

	distance := dotToAdd.distToCenter(bucket.Center)

	distsToCenter := SortedPairsByValue(bucket.DotsToDistToCenter, reverseLess[float64])

	if distance <= distsToCenter[0].Value {
		return true
	}

	grow := (distance - distsToCenter[0].Value) / (distsToCenter[0].Value - distsToCenter[1].Value)

	return grow < 2.0
}

func (bucket *Bucket) addDot(dot *Dot) {
	bucketDotsIds := keys(bucket.DotsToDistToCenter)
	bucketDots := getDots(bucketDotsIds)

	center := calculateCenter(bucketDots)
	bucket.Center = center

	bucket.DotsToDistToCenter[dot.ID] = dot.distToCenter(bucket.Center)
}

// ищем в бакете 2 выброса
// сначала находим те что максимально удалены от центра
// потом находим те у которых расстояние между ними больше 2*distanceToCenter
func (bucket *Bucket) findTwoAnomalyDots() (string, string) {
	distsToCenter := SortedPairsByValue(bucket.DotsToDistToCenter, reverseLess[float64])

	first := distsToCenter[0]
	firstDot := getDot(first.Key)

	for i := 1; i < len(distsToCenter); i++ {
		other := distsToCenter[i]
		otherDot := getDot(other.Key)
		scalar := scalarProduct(firstDot.History, otherDot.History)

		if scalar/(first.Value*other.Value) < 0 {
			return first.Key, other.Key
		}
		// ищем первый выброс с углом больше 90 градусов

	}

	panic("no anomaly dots found")
}

func (bucket *Bucket) oneAnomalyDot() string {
	distsToCenter := SortedPairsByValue(bucket.DotsToDistToCenter, reverseLess[float64])

	first := distsToCenter[0]
	return first.Key
}

// можем разбить бакет на 2 когда, 2 корректных выброса имеют diff больше maxAngle
// первый возвращаемое значение - точка по которой делим
// второй - true если можно делить, false если нельзя
// третий - все точки отсортированные по расстоянию до точки по которой делим
func (bucket *Bucket) canSplitBucketOver(firstDot *Dot) (bool, []Pair[string, float64], []Pair[string, float64]) {
	keys := keys(bucket.DotsToDistToCenter)
	dots := getDots(keys)

	firstsDists := make([]Pair[string, float64], 0, len(dots))

	for _, dot := range dots {
		firstsDists = append(firstsDists, Pair[string, float64]{Key: dot.ID, Value: distBetweenDots(firstDot, &dot)})
	}

	sort.Slice(firstsDists, func(i, j int) bool {
		return defaultLess(firstsDists[i].Value, firstsDists[j].Value)
	})

	index, maxGrows := maxGrow(firstsDists)

	return maxGrows >= 2.0, firstsDists[:index], firstsDists[index:]
}

func processSplitBucket(bucket *Bucket) (*Bucket, *Bucket) {
	anomaly := bucket.oneAnomalyDot()
	anomalyDot := getDot(anomaly)
	ok, firstDotPairs, secondDotPairs := bucket.canSplitBucketOver(&anomalyDot)
	if !ok {
		return nil, nil
	}

	firstDotsIds := make([]string, 0, len(firstDotPairs))
	secondDotsIds := make([]string, 0, len(secondDotPairs))

	for _, dot := range firstDotPairs {
		firstDotsIds = append(firstDotsIds, dot.Key)
	}

	for _, dot := range secondDotPairs {
		secondDotsIds = append(secondDotsIds, dot.Key)
	}

	firstDots := getDots(firstDotsIds)
	secondDots := getDots(secondDotsIds)

	firstBucket := proceessCreateBucket(firstDots)
	secondBucket := proceessCreateBucket(secondDots)

	return firstBucket, secondBucket
}

func proceessCreateBucket(dots []Dot) *Bucket {
	center := calculateCenter(dots)
	distsToCenter := make(map[string]float64)
	for _, dot := range dots {
		diff := diff(center, dot.History)
		distsToCenter[dot.ID] = dist(diff)
	}
	return &Bucket{
		Center:             center,
		DotsToDistToCenter: distsToCenter,
	}
}

func mergeBuckets(bucket1 *Bucket, bucket2 *Bucket) *Bucket {
	first := dirtyMultiply(bucket1.Center, float64(len(bucket1.DotsToDistToCenter)))
	second := dirtyMultiply(bucket2.Center, float64(len(bucket2.DotsToDistToCenter)))
	totalLen := float64(len(bucket1.DotsToDistToCenter) + len(bucket2.DotsToDistToCenter))
	center := dirtyDivide(plus(first, second), totalLen)

	distsToCenter := make(map[string]float64)
	for dotID, _ := range bucket1.DotsToDistToCenter {
		dot := getDot(dotID)
		distsToCenter[dotID] = dot.distToCenter(center)
	}

	for dotID, _ := range bucket2.DotsToDistToCenter {
		dot := getDot(dotID)
		distsToCenter[dotID] = dot.distToCenter(center)
	}

	return &Bucket{
		Center:             center,
		DotsToDistToCenter: distsToCenter,
	}
}
