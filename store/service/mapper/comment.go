package mapper

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
	"time"
)

func ToCommentDtos(comments []entity.Comment) []dto.CommentDto {
	dtoComments := make([]dto.CommentDto, 0, len(comments))
	for _, comment := range comments {
		dtoComments = append(dtoComments, *ToCommentDto(&comment))
	}
	return dtoComments
}

func ToCommentDto(comment *entity.Comment) *dto.CommentDto {
	return &dto.CommentDto{
		ID:        comment.ID,
		VideoID:   comment.VideoID,
		UpdatedAt: comment.UpdatedAt,
		Message:   comment.Message,
		Likes:     comment.Likes,
		Dislikes:  comment.Dislikes,
		Childs:    make([]dto.CommentDto, 0),
		ParentID:  comment.ParentID,
		RootID:    comment.RootID,
		Author:    *toCommentAuthorDto(&comment.Author),
	}
}

func ToComment(comment *dto.CommentDto) *entity.Comment {
	author := *toCommentAuthor(&comment.Author)
	commentEntity := &entity.Comment{
		ID:        comment.ID,
		VideoID:   comment.VideoID,
		Message:   comment.Message,
		Likes:     comment.Likes,
		Dislikes:  comment.Dislikes,
		Author:    author,
		ParentID:  comment.ParentID,
		RootID:    comment.RootID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	commentEntity.Relevance = calculateRelevance(commentEntity)

	return commentEntity
}

func toCommentAuthor(author *dto.CommentAuthorDto) *entity.CommentAuthor {
	return &entity.CommentAuthor{
		ID:       author.ID,
		Username: author.Username,
	}
}

func toCommentAuthorDto(author *entity.CommentAuthor) *dto.CommentAuthorDto {
	return &dto.CommentAuthorDto{
		ID:       author.ID,
		Username: author.Username,
	}
}

func calculateRelevance(comment *entity.Comment) float64 {
	return (float64(comment.Likes) - float64(comment.Dislikes)) + (30 - daysLater(comment.UpdatedAt))
}

func daysLater(createdAt time.Time) float64 {
	return time.Now().Sub(createdAt).Hours() / 24
}
