package service

import (
	"dev/bluebasooo/video-platform/api/dto"
)

const ChunckSz = 1024 * 1024

func GeneratePlan(meta *dto.PlainFileMetaDto) (*dto.WritePlanDto, error) {
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

func Write(taskID string, hash string, bytes []byte) error {
	cacheTasks.doneOp(taskID, hash)

	if cacheTasks.checkOfDone(taskID) {
		cacheTasks.doneTask(taskID)
	}

	// TODO: save bytes to file

	return nil
}
