package repo

import (
	obj_storage "dev/bluebasooo/video-platform/obj-storage"
	"dev/bluebasooo/video-platform/repo/entity"
	"fmt"
)

type FileRepository struct {
	objectStorage obj_storage.ObjectStorage
	bucketName    string
}

func NewFileRepository(objectStorage obj_storage.ObjectStorage) *FileRepository {
	return &FileRepository{objectStorage: objectStorage}
}

func (r *FileRepository) UploadFilePart(file *entity.FilePart) error {
	filePath := generateFilePath(file.FromUser, file.FileID, file.ID)
	return r.objectStorage.UploadFile(filePath, r.bucketName, file.Resource)
}

func (r *FileRepository) DownloadFilePart(file *entity.FilePart) ([]byte, error) {
	filePath := generateFilePath(file.FromUser, file.FileID, file.ID)
	return r.objectStorage.DownloadFile(filePath, r.bucketName)
}

func generateFilePath(userId string, fileID string, part string) string {
	return fmt.Sprintf("%s/%s/%s", userId, fileID, part)
}
