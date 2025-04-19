package repo

import (
	"context"
	"dev/bluebasooo/video-platform/db"
	"dev/bluebasooo/video-platform/repo/entity"

	"github.com/mongodb/mongo-go-driver/bson"
)

type FileMetaRepository struct {
	db *db.MongoDB
}

func NewFileMetaRepository(db *db.MongoDB) *FileMetaRepository {
	return &FileMetaRepository{db: db}
}

// fileId - приходит вместе в preview
func (r *FileMetaRepository) GetFileMeta(fileId string) (*entity.FileMeta, error) {
	collection := r.db.GetCollection("file_meta")

	filter := bson.M{"_id": fileId}
	var fileMeta entity.FileMeta
	err := collection.FindOne(context.TODO(), filter).Decode(&fileMeta)
	if err != nil {
		return nil, err
	}

	return &fileMeta, nil
}

func (r *FileMetaRepository) CreateFileMeta(fileMeta *entity.FileMeta) error {
	collection := r.db.GetCollection("file_meta")
	_, err := collection.InsertOne(context.TODO(), fileMeta)
	return err
}

// TODO update + delete
