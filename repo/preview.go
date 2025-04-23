package repo

import (
	"context"

	"dev/bluebasooo/video-platform/db"
	"dev/bluebasooo/video-platform/repo/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type PreviewRepository struct {
	db *db.MongoDB
}

func NewPreviewRepository(db *db.MongoDB) *PreviewRepository {
	return &PreviewRepository{db: db}
}

func (r *PreviewRepository) GetVideoPreview(videoID string) (*entity.VideoPreview, error) {
	collection := r.db.GetCollection("video_previews")

	filter := bson.M{"_id": videoID}
	var videoPreview entity.VideoPreview
	err := collection.FindOne(context.TODO(), filter).Decode(&videoPreview)
	if err != nil {
		return nil, err
	}

	return &videoPreview, nil
}

func (r *PreviewRepository) CreateVideoPreview(videoPreview *entity.VideoPreview) error {
	collection := r.db.GetCollection("video_previews")

	_, err := collection.InsertOne(context.TODO(), videoPreview)
	return err
}

func (r *PreviewRepository) GetVideoPreviews(ids []string) ([]entity.VideoPreview, error) {
	collection := r.db.GetCollection("video_previews")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	var videoPreviews []entity.VideoPreview
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &videoPreviews)
	if err != nil {
		return nil, err
	}

	return videoPreviews, nil
}
