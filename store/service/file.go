package service

import (
	"dev/bluebasooo/video-platform/repo"
	"dev/bluebasooo/video-platform/repo/entity"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	mockedFileChunks   map[string][]byte
	isMocked           bool                     = true
	fileMetaRepository *repo.FileMetaRepository // tmp need service type
)

func GetFilePartInterval(userId string, id string, part string) (map[string][]byte, error) {
	if isMocked {
		return MockedFileRead()
	}

	meta, err := fileMetaRepository.GetFileMeta(id)
	if err != nil {
		return nil, err
	}

	interval := getIntervalOfParts(meta.PartSequence, part)
	fileChunks := make(map[string][]byte)
	for _, part := range interval {
		filePart := entity.FilePart{
			FileID:   id,
			ID:       part,
			FromUser: userId,
		}
		partBytes, err := fileRepository.DownloadFilePart(&filePart)
		if err != nil {
			return nil, err
		}
		fileChunks[part] = partBytes
	}

	return fileChunks, nil
}

func getIntervalOfParts(parts []string, target string) []string {
	// at now make 1 interval before and 1 after

	// find target
	targetIndex := -1
	for i, part := range parts {
		if part == target {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return nil // not found
	}

	// very soft
	startInclusive := targetIndex - 1
	if startInclusive < 0 {
		startInclusive = 0
	}
	endExclusive := targetIndex + 1
	if endExclusive >= len(parts) {
		endExclusive = len(parts) - 1
	}

	// guaranted that slice does not changes after
	return parts[startInclusive:endExclusive]
}

func MockedFileRead() (map[string][]byte, error) {
	if mockedFileChunks != nil {
		return mockedFileChunks, nil
	}

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

	// Создаем WaitGroup для синхронизации горутин
	var wg sync.WaitGroup

	// Запускаем горутины для каждого чанка
	fileChanks := make(map[string][]byte)
	for i := range chunks {
		wg.Add(1) // Увеличиваем счетчик перед запуском новой горутины

		go func(chunk int) {
			defer wg.Done() // Уменьшаем счетчик после завершения горутины

			// Создаем буфер для чтения чанка
			buffer := make([]byte, chunkSize)

			// Читаем чанк из файла, начиная с позиции chunk * chunkSize
			_, err := file.ReadAt(buffer, int64(chunk)*int64(chunkSize))
			if err != nil {
				log.Printf("Error reading chunk %d: %v", chunk, err)
				return
			}

			// Здесь можно обрабатывать прочитанный чанк
			fileChanks[fmt.Sprintf("%d", chunk)] = buffer
		}(i)
	}

	wg.Wait() // Ждем завершения всех горутин

	mockedFileChunks = fileChanks
	return fileChanks, nil
}
