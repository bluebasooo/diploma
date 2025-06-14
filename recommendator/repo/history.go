package repo

import (
	"context"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/entity"
	"fmt"
	"strings"
)

type HistoryRepo struct {
	db *db.ClickhouseDB
}

func NewHistoryRepo(db *db.ClickhouseDB) *HistoryRepo {
	return &HistoryRepo{db: db}
}

const (
	getHistoryByUserIds = `
	SELECT user_id, video_id, sum(value) as metric 
	FROM history
	WHERE user_id IN (%s) 
		AND created_at >= now() - INTERVAL 30 DAY
	GROUP BY user_id, video_id
	`

	getHistoryAbout30Days = `
	SELECT user_id, created_at, video_id
	FROM history
	WHERE user_id = %s
		AND created_at >= now() - INTERVAL 30 DAY
	`

	insertHistory = `
	INSERT INTO history(user_id, video_id, created_at) VALUES 
	%s
	`
)

// Алгоритм
// 1. Получаем историю из метрик
// 2. Insert в историю
// 3. Для метрик по user пересчитываем историю

func (r *HistoryRepo) GetHistoryByUserIds(ctx context.Context, userIds []string) ([]entity.UserHistory, error) {
	rows, err := r.db.Db().Query(ctx, getHistoryByUserIds, userIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []entity.UserHistory
	for rows.Next() {
		var history entity.UserHistory
		err = rows.Scan(&history)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (r *HistoryRepo) GetHistoryAbout30Days(ctx context.Context, userId string) ([]entity.ShortUserHistory, error) {
	rows, err := r.db.Db().Query(ctx, getHistoryAbout30Days, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []entity.ShortUserHistory
	for rows.Next() {
		var history entity.ShortUserHistory
		err = rows.Scan(&history)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (r *HistoryRepo) BatchInsertHistory(ctx context.Context, histories []entity.History) error {
	values := make([]string, 0, len(histories))
	for _, history := range histories {
		values = append(values, fmt.Sprintf("'%s', '%s', '%s'", history.UserID, history.VideoID, history.CreatedAt))
	}
	valuesStr := strings.Join(values, ",\n")

	return r.db.Db().AsyncInsert(ctx, insertHistory, false, valuesStr)
}
