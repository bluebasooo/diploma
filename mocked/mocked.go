package mocked

import "dev/bluebasooo/video-platform/api/dto"

func MockedVideoPreview() *dto.VideoPreviewDto {
	return &dto.VideoPreviewDto{
		Name: "Mocked Video",
		Author: dto.AuthorShortPreviewDto{
			ID:   "1",
			Name: "Mocked Author",
			Img:  "https://localhost-img/150",
		},
		DurationMs: 10000,
		Stats: dto.VideoStatsDto{
			Views:      100,
			Likes:      100,
			Dislikes:   100,
			CommentsId: 100,
		},
	}
}
