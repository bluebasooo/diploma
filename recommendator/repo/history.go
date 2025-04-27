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
	getHistory = `
	SELECT %s 
	FROM history
	WHERE 1 = 1
	%s
	ORDER BY created_at DESC
	`

	insertHistory = `
	INSERT INTO history(%s)
	VALUES (
	%s
	)
	`
)

func (r *HistoryRepo) GetHistoryByUserId(ctx context.Context, userId string) ([]entity.History, error) {
	selector := "id, user_id, video_id, created_at"
	where := fmt.Sprintf("AND user_id = '%s' AND created_at >= now() - INTERVAL 30 DAY", userId)

	rows, err := r.db.Db().Query(ctx, getHistory, selector, where)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []entity.History
	for rows.Next() {
		var history entity.History
		err = rows.Scan(&history)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (r *HistoryRepo) BatchInsertHistory(ctx context.Context, histories []entity.History) error {
	plainColumns := "user_id, video_id, created_at"
	values := make([]string, 0, len(histories))
	for _, history := range histories {
		values = append(values, fmt.Sprintf("'%s', '%s', '%s'", history.UserID, history.VideoID, history.CreatedAt))
	}
	valuesStr := strings.Join(values, ",\n")

	return r.db.Db().AsyncInsert(ctx, insertHistory, false, plainColumns, valuesStr)
}
