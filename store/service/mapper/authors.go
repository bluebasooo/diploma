package mapper

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
)

func ToAuthorDto(author *entity.Author, subscribtions []entity.Author) *dto.AuthorDto {
	return &dto.AuthorDto{
		ID:       author.ID,
		Username: author.Username,
		ImgLink:  author.ImgLink,
		Stats: dto.StatsDto{
			Subscribers: author.AuthorStats.Subscribers,
			Likes:       author.AuthorStats.Likes,
			Views:       author.AuthorStats.Views,
		},
		Subscriptions: toSubscriptionDto(subscribtions),
	}
}

func toSubscriptionDto(authors []entity.Author) []dto.SubscriptionDto {
	subscriptions := make([]dto.SubscriptionDto, 0, len(authors))

	for _, author := range authors {
		subscriptions = append(subscriptions, dto.SubscriptionDto{
			ID:       author.ID,
			Username: author.Username,
			ImgLink:  author.ImgLink,
		})
	}

	return subscriptions
}
