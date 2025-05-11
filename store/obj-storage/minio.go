package obj_storage

import (
	"bytes"
	"context"
	"dev/bluebasooo/video-platform/config"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage struct {
	minioClient *minio.Client
	config      *config.MinioConfig
}

func NewObjectStorage(config *config.MinioConfig) *ObjectStorage {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: false,
	}
	minioClient, err := minio.New(config.Uri(), opts)
	if err != nil {
		log.Fatalf("Failed to create minio client: %v", err)
	}

	err = upsertBucket(minioClient, config.BucketName)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	return &ObjectStorage{minioClient: minioClient, config: config}
}

func upsertBucket(minioClient *minio.Client, bucketName string) error {
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ObjectStorage) UploadFile(objName string, bucketName string, paylaod []byte) error {
	reader := bytes.NewReader(paylaod)
	opts := minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	}
	_, err := s.minioClient.PutObject(context.TODO(), bucketName, objName, reader, int64(len(paylaod)), opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *ObjectStorage) DownloadFile(filePath string, bucketName string) ([]byte, error) {
	obj, err := s.minioClient.GetObject(context.TODO(), bucketName, filePath, minio.GetObjectOptions{})
	defer obj.Close()
	if err != nil {
		return nil, err
	}

	reader, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
