package service

import (
	"context"
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/service/mapper"
)

func GetUserHistory(userId string) ([]dto.HistoryDto, error) {
	userHistory, err := historyRepo.GetHistoryAbout30Days(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	result := make([]dto.HistoryDto, 0, len(userHistory))
	for _, val := range userHistory {
		result = append(result, *mapper.ToHistoryDto(&val))
	}

	return result, nil
}
