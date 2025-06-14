package mapper

import (
	"dev/bluebasooo/video-recomendator/api/dto"
)

func ToPoolDto(videoIds []string, page int, sz int) *dto.PoolDto {
	from := (page - 1) * sz
	to := page * sz

	return &dto.PoolDto{
		VideoIds: videoIds[from:to],
	}
}
