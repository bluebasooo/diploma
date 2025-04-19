package service

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	mockedFileChunks map[string][]byte
)

func GetFilePartInterval(id string, part string) (map[string][]byte, error) {
	return MockedFileRead()
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
