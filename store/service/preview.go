package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
	"dev/bluebasooo/video-platform/service/mapper"
)

func GetVideoPreview(id string) (*dto.VideoPreviewDto, error) {
	videoPreview, err := previewRepo.GetVideoPreview(id)
	if err != nil {
		return nil, err
	}

	authorEntity, err := authorRepository.GetAuthor(videoPreview.AuthorId) // TODO сделать заполнения авторских видео

	return mapper.ToVideoPreviewDto(videoPreview, authorEntity), nil
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

func CreateVideoPreview(id string, userId string, previewDto *dto.CreateVideoPreviewDto) error {
	videoPreview := mapper.ToVideoPreviewEntity(id, userId, previewDto)

	err := previewRepo.CreateVideoPreview(videoPreview)
	if err != nil {
		return err
	}

	go func() {
		author, err := authorRepository.GetAuthor(userId)
		if err != nil {
			return
		}
		bulk := make([]entity.VideoIndex, 0, 1)
		bulk = append(bulk, *mapper.ToVideoIndex(videoPreview, author.Username))
		Schedule(Index, "video-index", videoPreview.ID, bulk)
	}()

	return nil
}

func batchGetVideoPreviews(ids []string) ([]dto.VideoShortPreviewDto, error) {
	videoPreviews, err := previewRepo.GetVideoPreviews(ids)
	if err != nil {
		return nil, err
	}

	authorIds := make([]string, 0, len(videoPreviews))
	for _, videoPreview := range videoPreviews {
		authorIds = append(authorIds, videoPreview.AuthorId)
	}
	authorById, err := GetAuthorsUserNamesByIds(authorIds)
	if err != nil {
		return nil, err
	}

	videoShortPreviews := make([]dto.VideoShortPreviewDto, 0, len(videoPreviews))
	for _, videoPreview := range videoPreviews {
		author := authorById[videoPreview.AuthorId]
		videoShortPreviews = append(videoShortPreviews, *mapper.ToVideoShortPreviewDto(&videoPreview, &author))
	}

	return videoShortPreviews, nil
}
