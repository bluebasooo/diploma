package service

import (
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/engine"
	"dev/bluebasooo/video-recomendator/service/mapper"
)

func GetPagedVideoPool(id string, page int, pageSize int) *dto.PoolDto {
	pool := engine.GetPool(id)

	from := (page - 1) * pageSize
	to := page * pageSize

	return mapper.ToPoolDto(pool, from, to)
}
