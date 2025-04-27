package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/entity"
	"dev/bluebasooo/video-recomendator/repo"
	"dev/bluebasooo/video-recomendator/service/mapper"
)

var (
	historyRepo *repo.HistoryRepo
)

func GetHistoryByUserId(ctx context.Context, userId string) ([]dto.HistoryDto, error) {
	histories, err := historyRepo.GetHistoryByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	historyDtos := make([]dto.HistoryDto, 0, len(histories))
	for _, history := range histories {
		historyDtos = append(historyDtos, *mapper.ToHistoryDto(&history))
	}

	return historyDtos, nil
}

func BatchInsertHistory(ctx context.Context, historyDtos []dto.HistoryDto) error {
	histories := make([]entity.History, 0, len(historyDtos))
	for _, history := range historyDtos {
		histories = append(histories, *mapper.ToHistory(&history))
	}
	return historyRepo.BatchInsertHistory(ctx, histories)
}
