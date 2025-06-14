package mapper

import (
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/entity"
)

func ToHistoryDto(history *entity.ShortUserHistory) *dto.HistoryDto {
	return &dto.HistoryDto{
		UserID:    history.UserID,
		VideoID:   history.VideoID,
		CreatedAt: history.CreatedAt,
	}
}
