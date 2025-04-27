package mapper

import (
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/engine"
)

func ToPoolDto(pool *engine.Pool, start int, end int) *dto.PoolDto {
	pairs := engine.RangedKeysPairs(pool.RangedVideoIDs, start, end)

	return &dto.PoolDto{
		VideoIds: pairs,
	}
}
