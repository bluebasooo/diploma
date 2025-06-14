package repo

import (
	"context"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/entity"
	"fmt"
	"strings"
	"time"
)

type MetricsRepo struct {
	db *db.ClickhouseDB
}

func NewMetricsRepo(db *db.ClickhouseDB) *MetricsRepo {
	return &MetricsRepo{db: db}
}

const (
	getLastUncommitedMetrics = `
	SELECT view_id
	FROM (
		SELECT 
			view_id, user_id, video_id,
			argMax(type, created_at) AS latest_type
		FROM metrics
		GROUP BY view_id
	)
	WHERE latest_type = 'END'
	`

	insertMetrics = `
	INSERT INTO metrics(
		user_id, 
		video_id, 
		view_id, 
		type,
		value,
		created_at
	) VALUES (
		%s
	)
	`

	getCalculatedHistory = `
	SELECT 
		user_id, video_id, min(created_at),
		SUM(
			CASE 
				WHEN type = 'START' THEN 0
				WHEN type = 'LIKE' THEN 1.0 * value
				WHEN type = 'DISLIKE' THEN -1.0 * value
				WHEN type = 'WATCH_TIME' THEN 0.1 * value
				WHEN type = 'SHARE' THEN 1.0 * value
				WHEN type = 'END' THEN 0.0
			END
		) as metric
	FROM metrics
	WHERE (view_id, user_id, video_id) IN (%s)
	GROUP BY view_id, user_id, video_id
	`
)

func (r *MetricsRepo) CommitMetrics(ctx context.Context, viewIdentifiers []entity.ViewIdentifier) {
	commits := make([]entity.Metric, 0, len(viewIdentifiers))

	for _, view := range viewIdentifiers {
		commits = append(commits, entity.Metric{
			UserID:    view.UserID,
			VideoID:   view.VideoID,
			ViewID:    view.ViewID,
			Type:      string(entity.MetricTypeProcessed),
			Value:     0,
			CreatedAt: time.Now(),
		})
	}

	r.BatchInsertMetrics(ctx, commits)
}

func (r *MetricsRepo) BatchInsertMetrics(ctx context.Context, metrics []entity.Metric) error {
	values := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		pattern := "'%s', '%s', '%s', '%s', %f, '%s'"
		value := fmt.Sprintf(
			pattern,
			metric.UserID,
			metric.VideoID,
			metric.ViewID,
			metric.Type,
			metric.Value,
			metric.CreatedAt.Format(time.RFC3339),
		)
		values = append(values, value)
	}
	valuesStr := strings.Join(values, ",\n")

	return r.db.Db().AsyncInsert(ctx, insertMetrics, false, valuesStr)
}

func (r *MetricsRepo) GetCalculatedHistory(ctx context.Context, viewIdentifiers []entity.ViewIdentifier) ([]entity.History, error) {
	values := make([]string, 0, len(viewIdentifiers))
	for _, viewIdentifier := range viewIdentifiers {
		values = append(values, fmt.Sprintf("('%s', '%s', '%s')", viewIdentifier.ViewID, viewIdentifier.UserID, viewIdentifier.VideoID))
	}
	valuesStr := strings.Join(values, ",")

	rows, err := r.db.Db().Query(ctx, getCalculatedHistory, valuesStr)
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

func (r *MetricsRepo) GetLastUncommitedMetrics(ctx context.Context) ([]entity.ViewIdentifier, error) {
	rows, err := r.db.Db().Query(ctx, getLastUncommitedMetrics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viewIdentifiers []entity.ViewIdentifier
	for rows.Next() {
		var viewIdentifier entity.ViewIdentifier
		err = rows.Scan(&viewIdentifier)
		if err != nil {
			return nil, err
		}
		viewIdentifiers = append(viewIdentifiers, viewIdentifier)
	}
	return viewIdentifiers, nil
}
