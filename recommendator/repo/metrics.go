package repo

import (
	"context"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/entity"
	"fmt"
	"strings"
)

type MetricsRepo struct {
	db *db.ClickhouseDB
}

func NewMetricsRepo(db *db.ClickhouseDB) *MetricsRepo {
	return &MetricsRepo{db: db}
}

const (
	getMetrics = `
	SELECT %s 
	FROM metrics
	WHERE 1 = 1
	%s
	ORDER BY created_at DESC
	`

	insertMetrics = `
	INSERT INTO metrics(%s)
	VALUES (
	%s
	)
	`
)

func (r *MetricsRepo) BatchInsertMetrics(ctx context.Context, metrics []entity.Metric) error {
	plainColumns := "user_id, video_id, type, value, created_at"
	values := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		pattern := "'%s', '%s', '%s', %f, '%s'"
		value := fmt.Sprintf(pattern, metric.UserID, metric.VideoID, metric.Type, metric.Value, metric.CreatedAt)
		values = append(values, value)
	}
	valuesStr := strings.Join(values, ",\n")

	return r.db.Db().AsyncInsert(ctx, insertMetrics, false, plainColumns, valuesStr)
}

func (r *MetricsRepo) GetMetrics(ctx context.Context, userId, videoId string) ([]entity.Metric, error) {
	selector := "user_id, video_id, type, value, created_at"
	where := fmt.Sprintf("AND user_id = '%s' AND video_id = '%s'", userId, videoId)

	rows, err := r.db.Db().Query(ctx, getMetrics, selector, where)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []entity.Metric
	for rows.Next() {
		var metric entity.Metric
		err = rows.Scan(&metric)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
