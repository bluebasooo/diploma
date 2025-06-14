package repo

import (
	"context"
	mongoDb "dev/bluebasooo/video-common/db"
	"dev/bluebasooo/video-recomendator/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type DotsRepo struct {
	db *mongoDb.MongoDB
}

func NewDotsRepo(db *mongoDb.MongoDB) *DotsRepo {
	return &DotsRepo{db: db}
}

func (r *DotsRepo) GetDots(ctx context.Context, dotIds []string) ([]entity.DotHistory, error) {
	collection := r.db.GetCollection("dots")
	dots := make([]entity.DotHistory, 0)

	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": dotIds}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &dots)
	if err != nil {
		return nil, err
	}
	// обработать случай если не найдено

	return dots, nil
}

func (r *DotsRepo) GetDot(ctx context.Context, dotId string) (*entity.DotHistory, error) {
	collection := r.db.GetCollection("dots")
	var dot entity.DotHistory
	cursor, err := collection.Find(ctx, bson.M{"_id": dotId})
	if err != nil {
		return nil, err
	}

	err = cursor.Decode(&dot)
	if err != nil {
		return nil, err
	}

	return &dot, nil
}

func (r *DotsRepo) CreateDots(ctx context.Context, dots []entity.DotHistory) error {
	collection := r.db.GetCollection("dots")

	docs := make([]interface{}, 0, len(dots))
	for _, dot := range dots {
		docs = append(docs, dot)
	}

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}
