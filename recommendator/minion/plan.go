package minion

import (
	"context"
	"dev/bluebasooo/video-recomendator/entity"
	"time"
)

var batchSize int64 = 10_000
var forceUpdateOn = time.Duration(5 * time.Minute)

func (m *GrandMinion) Pull(ctx context.Context) ([]entity.History, error) {
	index, err := m.indexRepo.GetIndex(ctx)
	if err != nil {
		return nil, err
	}

	history, err := m.historyRepo.GetHistoryAfterID(ctx, index.ID, batchSize)
	if err != nil {
		return nil, err
	}

	return history, nil
}