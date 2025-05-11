package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
	"dev/bluebasooo/video-platform/service/mapper"
	"dev/bluebasooo/video-platform/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAuthor(author *dto.CreateAuthorDto) error {
	id := primitive.NewObjectID()
	authorEntity := entity.Author{
		ID:       &id,
		Username: author.Username,
	}
	return authorRepository.CreateAuthor(&authorEntity)
}

func GetAuthor(id string) (*dto.AuthorDto, error) {
	author, err := authorRepository.GetAuthor(id)
	if err != nil {
		return nil, err
	}

	subscribtions, err := authorRepository.GetAuthors(author.AuthorSubscriptions)
	if err != nil {
		return nil, err
	}

	return mapper.ToAuthorDto(author, subscribtions), nil
}

func GetAuthorsUserNamesByIds(ids []string) (map[string]entity.Author, error) {
	normalizedIds := utils.NormalizeIds(ids)

	authors, err := authorRepository.GetAuthors(normalizedIds)

	if err != nil {
		return nil, err
	}

	mappedAuthors := make(map[string]entity.Author, len(authors))
	for _, author := range authors {
		mappedAuthors[author.ID.Hex()] = author
	}

	return mappedAuthors, nil
}
