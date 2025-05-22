package minion

import (
	"context"
	"dev/bluebasooo/video-recomendator/entity"
	"time"
)

var batchSize = 1000
var forceUpdateOn = time.Duration(5 * time.Minute)

func (m *GrandMinion) Search(ctx context.Context, userId string) ([]entity.History, error) {
	index, err := m.indexRepo.GetIndex(ctx)
	if err != nil {
		return nil, err
	}

	history, err := m.historyRepo.GetHistoryAfterID(ctx, index.ID)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (m )