package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/entity"
	"dev/bluebasooo/video-recomendator/service/mapper"
)

func WriteMetrics(ctx context.Context, metrics []dto.MetricDto) error {
	metricEntities := make([]entity.Metric, 0, len(metrics))
	for _, metric := range metrics {
		metricEntities = append(metricEntities, *mapper.ToMetric(&metric))
	}
	err := metricsRepo.BatchInsertMetrics(ctx, metricEntities)
	if err != nil {
		return err
	}

	updatesHandler.Increment(len(metrics))

	return nil
}
