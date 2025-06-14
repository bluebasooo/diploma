package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/service/mapper"
	"maps"
)

func GetPagedVideoPool(id string, page int, pageSize int) (*dto.PoolDto, error) {
	dot, err := DotsRepo.GetDot(context.Background(), id)
	if err != nil {
		return nil, err
	}
	bucket, err := BucketRepo.GetBucket(context.Background(), dot.BucketID)
	if err != nil {
		return nil, err
	}

	center := maps.Clone(bucket.BucketCenter)
	videoIds := sortedByValueKeys(center, func(first float64, second float64) bool {
		return first > second
	})

	return mapper.ToPoolDto(videoIds, page, pageSize), nil
}
