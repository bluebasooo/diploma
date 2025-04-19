package dto

type WritePlanDto struct {
	TaskID string         `json:"taskId"`
	Ops    []OperationDto `json:"ops"`
}

type OperationDto struct {
	HashOperation string `json:"hashOperation"`
	BytesFrom     int64  `json:"bytesFrom"`
	BytesTo       int64  `json:"bytesTo"`
}
