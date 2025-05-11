package mapper

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
	"time"
)

func ToVideoPreviewEntity(id string, userId string, previewDto *dto.CreateVideoPreviewDto) *entity.VideoPreview {
	return &entity.VideoPreview{
		ID:          id,
		Name:        previewDto.Name,
		Img:         "", // TODO
		Description: previewDto.Description,
		AuthorId:    userId,
		VideoStats: entity.VideoStats{
			Views:    0,
			Likes:    0,
			Dislikes: 0,
		},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DurationMs: previewDto.DurationMs,
	}
}

func ToVideoPreviewDto(videoPreview *entity.VideoPreview, author *entity.Author) *dto.VideoPreviewDto {
	authorPreview := toAuthorShortPreviewDto(author)
	videoStats := toVideoStatsDto(&videoPreview.VideoStats)

	return &dto.VideoPreviewDto{
		Name:        videoPreview.Name,
		Description: videoPreview.Description,
		Author:      *authorPreview,
		DurationMs:  videoPreview.DurationMs,
		Stats:       *videoStats,
		CreatedAt:   videoPreview.CreatedAt,
	}
}

func toAuthorShortPreviewDto(author *entity.Author) *dto.AuthorShortPreviewDto {
	return &dto.AuthorShortPreviewDto{
		Username: author.Username,
		Img:      author.ImgLink,
		SubsNum:  author.AuthorStats.Subscribers,
	}
}

func toVideoStatsDto(videoStats *entity.VideoStats) *dto.VideoStatsDto {
	return &dto.VideoStatsDto{
		Views:    videoStats.Views,
		Likes:    videoStats.Likes,
		Dislikes: videoStats.Dislikes,
	}
}

func ToVideoShortPreviewDto(videoPreview *entity.VideoPreview, author *entity.Author) *dto.VideoShortPreviewDto {
	authorPreview := toAuthorShortPreviewDto(author)

	return &dto.VideoShortPreviewDto{
		ID:         videoPreview.ID,
		Img:        videoPreview.Img,
		Name:       videoPreview.Name,
		DurationMs: videoPreview.DurationMs,
		Author:     *authorPreview,
		Views:      videoPreview.VideoStats.Views,
	}
}
