package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo"
	"dev/bluebasooo/video-platform/repo/entity"
	"dev/bluebasooo/video-platform/service/mapper"
)

var authorRepository *repo.AuthorRepository

func CreateAuthor(author *dto.CreateAuthorDto) error {
	authorEntity := entity.Author{
		Username: author.Username,
	}
	return authorRepository.CreateAuthor(&authorEntity)
}

func GetAuthor(id string) (*dto.AuthorDto, error) {
	author, err := authorRepository.GetAuthor(id)
	if err != nil {
		return nil, err
	}

	subscribtions, err := authorRepository.GetAuthors(author.SubscriptionsOnAuthorsIDs)
	if err != nil {
		return nil, err
	}

	return mapper.ToAuthorDto(author, subscribtions), nil
}
