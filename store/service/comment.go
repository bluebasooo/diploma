package service

import (
	"context"
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
	"dev/bluebasooo/video-platform/service/mapper"
	"dev/bluebasooo/video-platform/utils"
)

func GetComments(ctx context.Context, videoId string, pageNum int, pageSize int) ([]dto.CommentDto, error) {
	roots, err := commRepo.GetCommentsPage(ctx, videoId, pageNum, pageSize)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(roots))
	for _, root := range roots {
		ids = append(ids, root.ID)
	}

	childs, err := commRepo.GetCommentChilds(ctx, ids)
	if err != nil {
		return nil, err
	}

	commentDtos := createCommentTree(roots, childs)

	return commentDtos, nil
}

func createCommentTree(roots []entity.Comment, childs []entity.Comment) []dto.CommentDto {
	mappedRoots := mapper.ToCommentDtos(roots)
	mappedChilds := mapper.ToCommentDtos(childs)

	rootsMap := utils.ToMapByUniqueField(mappedRoots, func(c *dto.CommentDto) string { return c.ID })
	childsMap := utils.ToMapByUniqueField(mappedChilds, func(c *dto.CommentDto) string { return c.ID })

	for _, v := range childs {
		parentID := v.ParentID
		parent, ok := rootsMap[parentID]
		if ok {
			child, ok := childsMap[v.ID]
			if ok {
				parent.Childs = append(parent.Childs, child)
			}
			continue
		}

		parent, ok = childsMap[parentID]
		if ok {
			child, ok := childsMap[v.ID]
			if ok {
				parent.Childs = append(parent.Childs, child)
			}
			continue
		}
	}

	return mappedRoots
}

func CreateComment(ctx context.Context, comment *dto.CommentDto) error {
	commentEntity := mapper.ToComment(comment)
	return commRepo.CreateComment(ctx, commentEntity)
}
