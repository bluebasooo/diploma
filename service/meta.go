package service

import (
	"dev/bluebasooo/video-platform/models"
	"fmt"
	"os"
)

const (
	CHUNK_SIZE = 1024 * 1024 // 1MB
)

// TODO: реализовать - пока что мок
func GetFileMeta(id string) (*models.FileMeta, error) {
	chunkSize := CHUNK_SIZE
	// Открываем файл
	file, err := os.Open("./mocked/mock.mp4")
	if err != nil {
		return nil, err
	}
	defer file.Close() // Гарантируем закрытие файла после завершения функции

	// Получаем информацию о файле (размер и т.д.)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Вычисляем общий размер файла и количество чанков
	fileSize := fileInfo.Size()
	chunks := int(fileSize / int64(chunkSize))
	partsIDs := make([]string, 0, chunks)
	for i := 0; i < chunks; i++ {
		partsIDs = append(partsIDs, fmt.Sprintf("%d", i))
	}

	return &models.FileMeta{
		Size:     fileSize,
		PartsIDs: partsIDs,
	}, nil
}
