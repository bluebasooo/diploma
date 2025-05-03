package repo

import (
	"context"
	"dev/bluebasooo/video-platform/db"
	"dev/bluebasooo/video-platform/repo/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepo struct {
	db *db.MongoDB
}

func NewCommentRepo(db *db.MongoDB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment *entity.Comment) error {
	collection := r.db.GetCollection("comments")
	_, err := collection.InsertOne(ctx, comment)
	return err
}

func (r *CommentRepo) GetCommentsPage(ctx context.Context, videoId string, pageNum int, pageSize int) ([]entity.Comment, error) {
	collection := r.db.GetCollection("comments")
	opts := nextPageOpts(pageNum, pageSize)
	cursor, err := collection.Find(ctx, bson.M{"videoId": videoId, "parentId": nil}, opts)
	if err != nil {
		return nil, err
	}

	var comments []entity.Comment
	err = cursor.All(ctx, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) GetCommentChilds(ctx context.Context, ids []string) ([]entity.Comment, error) {
	collection := r.db.GetCollection("comments")
	cursor, err := collection.Find(ctx, bson.M{"rootId": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	var comments []entity.Comment
	err = cursor.All(ctx, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func nextPageOpts(currentPageNum int, pageSize int) *options.FindOptions {
	skip := int64(currentPageNum * pageSize)
	limit := int64(currentPageNum * (pageSize + 1))
	return &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.M{"relevance": -1},
	}
}
