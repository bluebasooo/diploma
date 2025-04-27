package mapper

import (
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/entity"
)

func ToMetric(metricDto *dto.MetricDto) *entity.Metric {
	return &entity.Metric{
		UserID:    metricDto.UserID,
		VideoID:   metricDto.VideoID,
		Type:      metricDto.Type,
		Value:     metricDto.Value,
		CreatedAt: metricDto.CreatedAt,
	}
}
