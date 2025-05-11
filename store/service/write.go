package service

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/repo/entity"
	"dev/bluebasooo/video-platform/utils"
	"log"
	"time"
)

const ChunckSz = 1024 * 1024

// TMP solution need service type

func GeneratePlan(meta *dto.FileMetaPlanDto) (*dto.WritePlanDto, error) {
	id := utils.GetRandomHashId()
	chunks := meta.SizeInBytes / ChunckSz
	overlap := meta.SizeInBytes%ChunckSz == 0
	if !overlap {
		chunks++
	}

	plan := dto.WritePlanDto{
		TaskID: id,
		Ops:    make([]dto.OperationDto, 0, chunks),
	}
	hashes, err := cacheTasks.enqueue(plan.TaskID, chunks)
	if err != nil {
		return nil, err
	}
	for i := int64(0); i < chunks-1; i++ {
		plan.Ops = append(plan.Ops, dto.OperationDto{
			HashOperation: hashes[i],
			BytesFrom:     int64(i) * ChunckSz,
			BytesTo:       int64(i+1) * ChunckSz,
		})
	}

	plan.Ops = append(plan.Ops, dto.OperationDto{
		HashOperation: hashes[chunks-1],
		BytesFrom:     int64(chunks-1) * ChunckSz,
		BytesTo:       meta.SizeInBytes,
	})

	en := entity.FileMeta{
		ID:           id,
		Parts:        make([]entity.FileMetaPart, 0, len(plan.Ops)),
		FullSz:       meta.SizeInBytes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PartSequence: make([]string, 0, len(plan.Ops)),
		IsDraft:      true,
	}
	for _, op := range plan.Ops {
		en.Parts = append(en.Parts, entity.FileMetaPart{
			Hash: op.HashOperation,
			Sz:   op.BytesTo - op.BytesFrom,
		})
		en.PartSequence = append(en.PartSequence, op.HashOperation)
	}
	fileMetaRepository.CreateFileMeta(&en)

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

	err := fileRepository.UploadFilePart(&filePart)
	if err != nil {
		cacheTasks.err(taskID, hash)
		log.Println("Failed to upload file part:", err)
		return err
	}

	cacheTasks.doneOp(taskID, hash)

	if cacheTasks.checkOfDone(taskID) {
		err = fileMetaRepository.CommitDraft(taskID)
		if err != nil {
			return err
		}
		cacheTasks.doneTask(taskID)
	}

	return nil
}
