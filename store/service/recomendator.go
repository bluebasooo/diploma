package service

import (
	"dev/bluebasooo/video-platform/mocked"
	"errors"
)

func GetRecommendedVideosIds(userID string) ([]string, error) {
	if mocked.IsMocked {
		return mocked.MockedRecommendedVideosIds(), nil
	}

	return nil, errors.New("not implemented")
}
