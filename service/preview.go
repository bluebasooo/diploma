package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/mocked"
)

func GetVideoPreview(id string) (*dto.VideoPreviewDto, error) {
	return mocked.MockedVideoPreview(), nil
}
