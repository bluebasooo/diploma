package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/service/mapper"
)

const (
	CHUNK_SIZE = 1024 * 1024 // 1MB
)

func GetFileMeta(id string) (*dto.FileMetaDto, error) {
	meta, err := fileMetaRepository.GetFileMeta(id)
	if err != nil {
		return nil, err
	}

	return mapper.ToFileMetaDto(meta), nil
}
