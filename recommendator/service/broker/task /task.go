package task

type TaskType = string

const (
	HistoryUpdate   TaskType = "history_update"
	UserUpdate      TaskType = "user_update"
	UpdateDots      TaskType = "update_dots"
	SeparateBuckets TaskType = "separate_buckets"
)
