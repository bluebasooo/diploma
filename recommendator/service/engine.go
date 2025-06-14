package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/entity"
	"dev/bluebasooo/video-recomendator/handler"
	"dev/bluebasooo/video-recomendator/repo"
	"dev/bluebasooo/video-recomendator/service/broker"
	"dev/bluebasooo/video-recomendator/utils"
	"log"
	"sync"
	"time"
)

var (
	historyRepo         *repo.HistoryRepo
	metricsRepo         *repo.MetricsRepo
	dotsRepo            *repo.DotsRepo
	bucketRepo          *repo.BucketRepo
	historyUpdateBroker = broker.NewDummyBroker(func(payload *HistoryUpdatesPayload) error {
		return HistoryUpdates(context.Background(), payload.ViewIdentifiers)
	})
	userUpdatesBroker = broker.NewDummyBroker(func(payload *UserUpdatesPayload) error {
		return UserUpdates(context.Background(), payload.UserIds)
	})
	separateBucketBroker = broker.NewDummyBroker(func(payload *SeparateBucketPayload) error {
		return SeparateBuckets(payload.bucketIds)
	})
	updatesHandler *handler.UpdateHandler
)

// Алгоритм
// 1. Берем все незакомиченные метрики - getLastUncomitedMetrics - view_id, user_id, video_id
// 2. По ним вычисляем агрегат - метрику просмотра + коммит просмотра
// 3. Планируем обновление - где пересчитываем точки по истории
// 4. Коммитим метрику - добавляем тип PROCESSED

func Loop() {
	errs := make(chan error)
	for {
		select {
		case <-updatesHandler.Producer:
			err := ProcessMetrics(context.Background())
			if err != nil {
				errs <- err
			}
		case err := <-errs:
			log.Print(err)
		}
	}
}

func ProcessMetrics(ctx context.Context) error {
	// 1.
	viewIdentifiers, err := metricsRepo.GetLastUncommitedMetrics(ctx)
	if err != nil {
		return err
	}

	err = historyUpdateBroker.AsyncExecution(&HistoryUpdatesPayload{
		ViewIdentifiers: viewIdentifiers,
	})
	if err != nil {
		return err
	}

	// 4.
	metricsRepo.CommitMetrics(ctx, viewIdentifiers)

	return nil
}

type HistoryUpdatesPayload struct {
	ViewIdentifiers []entity.ViewIdentifier
}

func HistoryUpdates(ctx context.Context, viewIdentifiers []entity.ViewIdentifier) error {
	history, err := metricsRepo.GetCalculatedHistory(ctx, viewIdentifiers)
	if err != nil {
		return err
	}

	err = historyRepo.BatchInsertHistory(ctx, history)
	if err != nil {
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
		return err
	}

	return nil
}

type UserUpdatesPayload struct {
	UserIds []string
}

func UserUpdates(ctx context.Context, userIds []string) error {
	userHistory, err := historyRepo.GetHistoryByUserIds(ctx, userIds)
	if err != nil {
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
		return err
	}

	err = dotsRepo.CreateDots(ctx, result.success)
	if err != nil {
		return err
	}
	err = bucketRepo.UpsertBuckets(ctx, result.bucketUpdates...)
	if err != nil {
		return err
	}

	bucketToSeparate := Plain(
		result.separation,
		func(bucketId string, separateVal int) string {
			return bucketId
		})

	err = separateBucketBroker.AsyncExecution(&SeparateBucketPayload{
		bucketIds: bucketToSeparate,
	})
	if err != nil {
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
	buckets, err := bucketRepo.GetAllBuckets(ctx)
	if err != nil {
		return nil, err
	}

	fail := make([]entity.DotHistory, 0)
	success := make([]entity.DotHistory, 0)
	bucketUpdates := make([]entity.Bucket, 0)

	trySeparate := make(map[string]int)
	for i, dot := range dots {

		var minimus float64
		var bestMatchBucketId string
		for bucketId, bucket := range buckets {
			dist, ok := ShouldAddToBucket(&bucket, &dot)
			if ok && minimus > dist {
				minimus = dist
				bestMatchBucketId = bucketId
			}
		}

		if bestMatchBucketId != "" {
			bucketToAdd := buckets[bestMatchBucketId]
			resultBucket, err := AddDot(&bucketToAdd, &dot)
			if err != nil {
				log.Print(err.Error())
				fail = append(fail, dot)
				continue
			}
			bucketUpdates = append(bucketUpdates, *resultBucket)

			dots[i].BucketID = bucketToAdd.ID
			success = append(success, dots[i])

			trySeparate[bestMatchBucketId] += 1
		}

	}

	return &AddingDotsResult{
		success:       success,
		fail:          fail,
		bucketUpdates: bucketUpdates,
		separation:    trySeparate,
	}, nil
}

type SeparateBucketPayload struct {
	bucketIds []string
}

func SeparateBuckets(bucketIds []string) error {
	buckets, err := bucketRepo.GetBuckets(context.Background(), bucketIds...)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	wg := sync.WaitGroup{}
	result := make(chan *entity.Bucket)
	defer close(result)

	for _, val := range buckets {
		wg.Add(1)
		go func(bucket *entity.Bucket) {
			defer wg.Done()

			oBucket, tBucket := ProcessSplitBucket(bucket)
			result <- oBucket
			result <- tBucket
		}(&val)
	}

	wg.Wait()

	updates := make([]entity.Bucket, 0, len(result))
	for bucket := range result {
		updates = append(updates, *bucket)
	}

	err = bucketRepo.UpsertBuckets(context.Background(), updates...)
	if err != nil {
		return err
	}

	return nil
}
