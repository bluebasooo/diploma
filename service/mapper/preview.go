package mapper

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
)

func ToVideoPreviewDto(videoPreview *entity.VideoPreview) *dto.VideoPreviewDto {
	authorPreview := toAuthorShortPreviewDto(&videoPreview.AuthorPreview)
	videoStats := toVideoStatsDto(&videoPreview.VideoStats)

	return &dto.VideoPreviewDto{
		Name:          videoPreview.Name,
		Description:   videoPreview.Description,
		Author:        *authorPreview,
		DurationMs:    videoPreview.DurationMs,
		Stats:         *videoStats,
		CreatedAt:     videoPreview.CreatedAt,
		CommentTreeID: videoPreview.CommentTreeId,
	}
}

func toAuthorShortPreviewDto(authorPreview *entity.AuthorPreview) *dto.AuthorShortPreviewDto {
	return &dto.AuthorShortPreviewDto{
		Name:    authorPreview.Name,
		Img:     authorPreview.Img,
		SubsNum: authorPreview.SubsNum,
	}
}

func toVideoStatsDto(videoStats *entity.VideoStats) *dto.VideoStatsDto {
	return &dto.VideoStatsDto{
		Views:    videoStats.Views,
		Likes:    videoStats.Likes,
		Dislikes: videoStats.Dislikes,
	}
}

func ToVideoShortPreviewDto(videoPreview *entity.VideoPreview) *dto.VideoShortPreviewDto {
	authorPreview := toAuthorShortPreviewDto(&videoPreview.AuthorPreview)

	return &dto.VideoShortPreviewDto{
		ID:         videoPreview.ID,
		Img:        videoPreview.Img,
		Name:       videoPreview.Name,
		DurationMs: videoPreview.DurationMs,
		Author:     *authorPreview,
		Views:      videoPreview.VideoStats.Views,
	}
}
