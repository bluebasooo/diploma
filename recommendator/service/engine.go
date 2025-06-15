package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/entity"
	"dev/bluebasooo/video-recomendator/handler"
	"dev/bluebasooo/video-recomendator/repo"
	"dev/bluebasooo/video-recomendator/service/broker"
	"dev/bluebasooo/video-recomendator/utils"
	"log"
	"math"
	"sync"
	"time"
)

var (
	HistoryRepo         *repo.HistoryRepo
	MetricsRepo         *repo.MetricsRepo
	DotsRepo            *repo.DotsRepo
	BucketRepo          *repo.BucketRepo
	historyUpdateBroker = broker.NewDummyBroker(func(payload *HistoryUpdatesPayload) error {
		return HistoryUpdates(context.Background(), payload.ViewIdentifiers)
	})
	userUpdatesBroker = broker.NewDummyBroker(func(payload *UserUpdatesPayload) error {
		return UserUpdates(context.Background(), payload.UserIds)
	})
	separateBucketBroker = broker.NewDummyBroker(func(payload *SeparateBucketPayload) error {
		return SeparateBuckets(payload.bucketIds)
	})
	UpdatesHandler *handler.UpdateHandler
)

// Алгоритм
// 1. Берем все незакомиченные метрики - getLastUncomitedMetrics - view_id, user_id, video_id
// 2. По ним вычисляем агрегат - метрику просмотра + коммит просмотра
// 3. Планируем обновление - где пересчитываем точки по истории
// 4. Коммитим метрику - добавляем тип PROCESSED

func Loop() {
	errs := make(chan error, 1000)

	go func() {
		historyUpdateBroker.EventLoop(errs)
	}()

	go func() {
		userUpdatesBroker.EventLoop(errs)
	}()

	go func() {
		separateBucketBroker.EventLoop(errs)
	}()

	for {
		select {
		case <-UpdatesHandler.Producer:
			err := ProcessMetrics(context.Background())
			if err != nil {
				errs <- err
			}
		case err := <-errs:
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func ProcessMetrics(ctx context.Context) error {
	// 1.
	viewIdentifiers, err := MetricsRepo.GetLastUncommitedMetrics(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if viewIdentifiers == nil {
		return nil
	}

	err = historyUpdateBroker.AsyncExecution(&HistoryUpdatesPayload{
		ViewIdentifiers: viewIdentifiers,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 4.
	MetricsRepo.CommitMetrics(ctx, viewIdentifiers)

	return nil
}

type HistoryUpdatesPayload struct {
	ViewIdentifiers []entity.ViewIdentifier
}

func HistoryUpdates(ctx context.Context, viewIdentifiers []entity.ViewIdentifier) error {
	history, err := MetricsRepo.GetCalculatedHistory(ctx, viewIdentifiers)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = HistoryRepo.BatchInsertHistory(ctx, history)
	if err != nil {
		log.Fatal(err)
		return err
	}

	userIdsToUpdate := utils.NewSet[string]()
	for _, view := range history {
		userIdsToUpdate.Add(view.UserID)
	}
	updates := userIdsToUpdate.AsArr()
	//UserUpdates(ctx, updates)
	err = userUpdatesBroker.AsyncExecution(&UserUpdatesPayload{
		UserIds: updates,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

type UserUpdatesPayload struct {
	UserIds []string
}

func UserUpdates(ctx context.Context, userIds []string) error {
	userHistory, err := HistoryRepo.GetHistoryByUserIds(ctx, userIds)
	if err != nil {
		log.Fatal(err)
		return err
	}

	userByHistory := GroupByValueProp(userHistory, func(u entity.UserHistory) string { return u.UserID })

	dots := make([]entity.DotHistory, 0)
	for userID, history := range userByHistory {
		historyMap := make(map[string]float64)
		for _, h := range history {
			historyMap[h.VideoID] += h.Metric
		}

		dot := entity.DotHistory{
			ID: entity.DotHistoryVersionedId{
				DotID:      userID,
				DateUpdate: time.Now(),
			},
			History: historyMap,
		}
		dots = append(dots, dot)
	}

	result, err := AddDotsToBucket(ctx, dots)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = DotsRepo.CreateDots(ctx, result.success)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//err = BucketRepo.UpsertBuckets(ctx, result.bucketUpdates...)
	//if err != nil {
	//	log.Fatal(err)
	//	return err
	//}

	bucketToSeparate := Plain(
		result.separation,
		func(bucketId string, separateVal int) string {
			return bucketId
		})

	err = separateBucketBroker.AsyncExecution(&SeparateBucketPayload{
		bucketIds: bucketToSeparate,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

type AddingDotsResult struct {
	success       []entity.DotHistory
	fail          []entity.DotHistory
	bucketUpdates []entity.Bucket
	separation    map[string]int
}

// try to update dot and separate bucket
func AddDotsToBucket(ctx context.Context, dots []entity.DotHistory) (*AddingDotsResult, error) {
	fail := make([]entity.DotHistory, 0)
	success := make([]entity.DotHistory, 0)
	//bucketUpdates := make([]entity.Bucket, 0)

	trySeparate := make(map[string]int)
	for i, dot := range dots {
		maybeExist, err := DotsRepo.GetDot(ctx, dot.GetDotID())
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		buckets, err := BucketRepo.GetAllBuckets(ctx)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		var bestMatchBucketId string
		if maybeExist == nil { // try to find best bucket for new dot
			var minimus = math.MaxFloat64
			var notFoundBucketId string
			for bucketId, bucket := range buckets {
				dist, ok := ShouldAddToBucket(&bucket, &dot)
				if ok && minimus >= dist {
					minimus = dist
					bestMatchBucketId = bucketId
				} else if minimus >= dist {
					notFoundBucketId = bucketId
				}
			}

			if bestMatchBucketId == "" {
				bestMatchBucketId = notFoundBucketId
			}
		} else {
			bestMatchBucketId = maybeExist.BucketID
		}

		bucketToAdd := buckets[bestMatchBucketId]
		resultBucket, err := AddDot(bucketToAdd.ID, &dot)
		if err != nil {
			log.Fatal(err)
			log.Print(err.Error())
			fail = append(fail, dot)
			continue
		}
		//bucketUpdates = append(bucketUpdates, *resultBucket)

		dots[i].BucketID = bucketToAdd.ID
		success = append(success, dots[i])

		if len(resultBucket.BucketDotsToDistToCenter) < 5 {
			continue
		}
		trySeparate[bestMatchBucketId] += 1

	}

	return &AddingDotsResult{
		success: success,
		fail:    fail,
		//bucketUpdates: bucketUpdates,
		separation: trySeparate,
	}, nil
}

type SeparateBucketPayload struct {
	bucketIds []string
}

func SeparateBuckets(bucketIds []string) error {
	wg := sync.WaitGroup{}

	for _, val := range bucketIds {
		wg.Add(1)
		go func(bucketId string) {
			defer wg.Done()
			var bucket *entity.Bucket
			for {
				if BucketRepo.LockOnBucket(val) {
					bucket, _ = BucketRepo.GetBucket(context.Background(), bucketId)
					break
				}
			}
			defer BucketRepo.UnlockOnBucket(bucketId)
			if bucket.IsSeparated {
				return
			}

			oBucket, tBucket := ProcessSplitBucket(bucket)
			if oBucket == nil && tBucket == nil {
				return
			}

			bucket.IsSeparated = true
			err := BucketRepo.UpsertBuckets(context.Background(), *bucket, *oBucket, *tBucket)
			if err != nil {
				log.Fatal(err)
			}
		}(val)
	}

	wg.Wait()

	return nil
}
