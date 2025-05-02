package mapper

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
)

func ToFileMetaDto(meta *entity.FileMeta) *dto.FileMetaDto {
	parts := make([]dto.FileMetaPartDto, 0, len(meta.Parts))
	for _, part := range meta.Parts {
		partDto := *toFileMetaPartDto(&part)
		parts = append(parts, partDto)
	}

	return &dto.FileMetaDto{
		ID:     meta.ID,
		Parts:  parts,
		FullSz: meta.FullSz,
	}
}

func toFileMetaPartDto(part *entity.FileMetaPart) *dto.FileMetaPartDto {
	return &dto.FileMetaPartDto{
		Hash:  part.Hash,
		Sz:    part.Sz,
		S3Url: part.S3Url,
	}
}
