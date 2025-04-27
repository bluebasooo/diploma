package mapper

import (
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/entity"
)

func ToHistory(history *dto.HistoryDto) *entity.History {
	return &entity.History{
		UserID:    history.UserID,
		VideoID:   history.VideoID,
		CreatedAt: history.CreatedAt,
		Metric:    history.Metric,
	}
}

func ToHistoryDto(history *entity.History) *dto.HistoryDto {
	return &dto.HistoryDto{
		UserID:    history.UserID,
		VideoID:   history.VideoID,
		CreatedAt: history.CreatedAt,
		Metric:    history.Metric,
	}
}
