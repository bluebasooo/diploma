package repo

import (
	"context"
	mongoDb "dev/bluebasooo/video-common/db"
	"dev/bluebasooo/video-recomendator/entity"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DotsRepo struct {
	db *mongoDb.MongoDB
}

func NewDotsRepo(db *mongoDb.MongoDB) *DotsRepo {
	return &DotsRepo{db: db}
}

func (r *DotsRepo) GetDots(ctx context.Context, dotIds []string) ([]entity.DotHistory, error) {
	collection := r.db.GetCollection("dots")
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id.id": bson.M{"$in": dotIds}}}},
		{{Key: "$sort", Value: bson.D{{"_id.id", 1}, {"_id.date_update", -1}}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{"id": "$_id.id"},
			"doc": bson.M{"$first": "$$ROOT"},
		}}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$doc"}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []entity.DotHistory
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *DotsRepo) GetDot(ctx context.Context, dotId string) (*entity.DotHistory, error) {
	collection := r.db.GetCollection("dots")

	filter := bson.M{"_id.id": dotId}
	opts := options.FindOne().SetSort(bson.D{{"_id.date_update", -1}})

	var result entity.DotHistory
	err := collection.FindOne(ctx, filter, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
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
