package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo"
)

var searchRepo *repo.SearchRepo

func FindVideos(query string) ([]dto.VideoShortPreviewDto, error) {
	res, err := searchRepo.SearchVideos(query)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(res))
	for _, video := range res {
		ids = append(ids, video.ID)
	}

	videoPreviews, err := batchGetVideoPreviews(ids)
	if err != nil {
		return nil, err
	}

	return videoPreviews, nil
}
