package mocked

import (
	"dev/bluebasooo/video-platform/api/dto"
	"time"
)

var (
	IsMocked = true
)

func MockedVideoPreview() *dto.VideoPreviewDto {
	return &dto.VideoPreviewDto{
		Name: "Mocked Video",
		Author: dto.AuthorShortPreviewDto{
			Name:    "Mocked Author",
			Img:     "https://localhost:8080/img/150",
			SubsNum: 100,
		},
		DurationMs: 10000,
		Stats: dto.VideoStatsDto{
			Views:    100,
			Likes:    100,
			Dislikes: 100,
		},
		CreatedAt:     time.Now(),
		CommentTreeID: "1",
	}
}

func MockedRecommendedVideosIds() []string {
	return []string{
		"1",
		"2",
		"3",
	}
}
