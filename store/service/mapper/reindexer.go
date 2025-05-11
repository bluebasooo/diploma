package mapper

import "dev/bluebasooo/video-platform/repo/entity"

func ToVideoIndex(videoPreview *entity.VideoPreview, authorName string) *entity.VideoIndex {
	return &entity.VideoIndex{
		Title:       videoPreview.Name,
		Description: videoPreview.Description,
		DurationMs:  int(videoPreview.DurationMs),
		AuthorName:  authorName,
		UploadDate:  videoPreview.CreatedAt,
		Views:       0,
		Hidden:      false,
	}
}
