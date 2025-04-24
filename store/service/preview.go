package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/mocked"
	"dev/bluebasooo/video-platform/repo"
	"dev/bluebasooo/video-platform/service/mapper"
)

var (
	previewRepo *repo.PreviewRepository
)

func GetVideoPreview(id string) (*dto.VideoPreviewDto, error) {
	if mocked.IsMocked {
		return mocked.MockedVideoPreview(), nil
	}

	videoPreview, err := previewRepo.GetVideoPreview(id)
	if err != nil {
		return nil, err
	}

	return mapper.ToVideoPreviewDto(videoPreview), nil
}

func GetMainPageVideoPreviews(userID string) ([]dto.VideoShortPreviewDto, error) {
	videosIds, err := GetRecommendedVideosIds(userID)
	if err != nil {
		return nil, err
	}

	videoPreviews, err := batchGetVideoPreviews(videosIds)
	if err != nil {
		return nil, err
	}

	return videoPreviews, nil
}

func batchGetVideoPreviews(ids []string) ([]dto.VideoShortPreviewDto, error) {
	videoPreviews, err := previewRepo.GetVideoPreviews(ids)
	if err != nil {
		return nil, err
	}

	videoShortPreviews := make([]dto.VideoShortPreviewDto, 0, len(videoPreviews))
	for _, videoPreview := range videoPreviews {
		videoShortPreviews = append(videoShortPreviews, *mapper.ToVideoShortPreviewDto(&videoPreview))
	}

	return videoShortPreviews, nil
}
