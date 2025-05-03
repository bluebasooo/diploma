package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo"
	"dev/bluebasooo/video-platform/repo/entity"
	"log"
)

const ChunckSz = 1024 * 1024

// TMP solution need service type
var fileRepository *repo.FileRepository

func GeneratePlan(meta *dto.FileMetaPlanDto) (*dto.WritePlanDto, error) {
	chunks := meta.SizeInBytes / ChunckSz
	overlap := meta.SizeInBytes%ChunckSz == 0
	if overlap {
		chunks++
	}

	plan := dto.WritePlanDto{
		TaskID: meta.TaskID,
		Ops:    make([]dto.OperationDto, 0, chunks),
	}
	hashes, err := cacheTasks.enqueue(plan.TaskID, chunks)
	if err != nil {
		return nil, err
	}
	for i := int64(0); i < chunks; i++ {
		plan.Ops = append(plan.Ops, dto.OperationDto{
			HashOperation: hashes[i],
			BytesFrom:     int64(i) * ChunckSz,
			BytesTo:       int64(i+1) * ChunckSz,
		})
	}

	return &plan, nil
}

func Write(taskID string, hash string, bytes []byte, userID string) error {
	filePart := entity.FilePart{
		FileID:   taskID, // idk - this file ID i think
		ID:       hash,   // from plan
		Sz:       int64(len(bytes)),
		Resource: bytes,
		FromUser: userID,
	}

	// io operation
	// more smart way on wait group
	go func() {
		err := fileRepository.UploadFilePart(&filePart)
		if err != nil {
			cacheTasks.err(taskID, hash)
			log.Println("Failed to upload file part:", err)
		}

		cacheTasks.doneOp(taskID, hash)

		if cacheTasks.checkOfDone(taskID) {
			cacheTasks.doneTask(taskID)
		}
	}()

	return nil
}
