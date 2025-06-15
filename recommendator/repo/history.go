package repo

import (
	"context"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/entity"
	"fmt"
	"strings"
	"time"
)

type HistoryRepo struct {
	db *db.ClickhouseDB
}

func NewHistoryRepo(db *db.ClickhouseDB) *HistoryRepo {
	return &HistoryRepo{db: db}
}

const (
	getHistoryByUserIds = `
	SELECT user_id, video_id, sum(value) as value 
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
	INSERT INTO history(user_id, video_id, created_at, value) VALUES 
	%s
	`
)

// Алгоритм
// 1. Получаем историю из метрик
// 2. Insert в историю
// 3. Для метрик по user пересчитываем историю

func (r *HistoryRepo) GetHistoryByUserIds(ctx context.Context, userIds []string) ([]entity.UserHistory, error) {
	pattern := "'%s'"
	patterned := make([]string, 0, len(userIds))
	for _, id := range userIds {
		patterned = append(patterned, fmt.Sprintf(pattern, id))
	}
	joined := strings.Join(patterned, ", ")
	query := fmt.Sprintf(getHistoryByUserIds, joined)
	rows, err := r.db.Db().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []entity.UserHistory
	for rows.Next() {
		var history entity.UserHistory
		err = rows.ScanStruct(&history)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (r *HistoryRepo) GetHistoryAbout30Days(ctx context.Context, userId string) ([]entity.ShortUserHistory, error) {
	query := fmt.Sprintf(getHistoryAbout30Days, userId)
	rows, err := r.db.Db().Query(ctx, query)
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
		values = append(values, fmt.Sprintf(
			"('%s', '%s', '%s', %f)",
			history.UserID,
			history.VideoID,
			history.CreatedAt.Format(time.DateTime),
			history.Metric,
		))
	}
	valuesStr := strings.Join(values, ",\n")
	insertQuery := fmt.Sprintf(insertHistory, valuesStr)

	return r.db.Db().AsyncInsert(ctx, insertQuery, true)
}
