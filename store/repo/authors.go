package repo

import (
	"context"
	"dev/bluebasooo/video-platform/db"
	"dev/bluebasooo/video-platform/repo/entity"

	"go.mongodb.org/mongo-driver/bson"
)

type AuthorRepository struct {
	db *db.MongoDB
}

func NewAuthorRepository(db *db.MongoDB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) CreateAuthor(author *entity.Author) error {
	_, err := r.db.GetCollection("authors").InsertOne(context.Background(), author)
	return err
}

func (r *AuthorRepository) GetAuthor(id string) (*entity.Author, error) {
	var author entity.Author
	err := r.db.GetCollection("authors").FindOne(context.Background(), bson.M{"_id": id}).Decode(&author)
	return &author, err
}

func (r *AuthorRepository) GetAuthors(ids []string) ([]entity.Author, error) {
	var authors []entity.Author
	cursor, err := r.db.GetCollection("authors").Find(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &authors)
	return authors, err
}
