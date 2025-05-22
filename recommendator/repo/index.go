package repo

import (
	"context"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/entity"
	"errors"
)

type IndexRepo struct {
	db *db.ClickhouseDB
}

func NewIndexRepo(db *db.ClickhouseDB) *IndexRepo {
	return &IndexRepo{db: db}
}

const (
	getIndex = `
	SELECT id, last_updated
	FROM index
	ORDER BY last_updated DESC
	LIMIT 1
	`

	insertIndex = `
	INSERT INTO index (id, last_updated)
	VALUES ($1, $2)
	`
)

func (r *IndexRepo) GetIndex(ctx context.Context) (entity.Index, error) {
	rows, err := r.db.Db().Query(ctx, getIndex)
	if err != nil {
		return entity.Index{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return entity.Index{}, errors.New("index not found")
	}

	var index entity.Index
	err = rows.Scan(&index.ID, &index.LastUpdated)
	if err != nil {
		return entity.Index{}, err
	}

	return index, nil
}

func (r *IndexRepo) CommitIndex(ctx context.Context, index entity.Index) error {
	err := r.db.Db().Exec(ctx, insertIndex, index.ID, index.LastUpdated)
	if err != nil {
		return err
	}

	return nil
}
