package entity

// storage
type FilePart struct {
	FileID   string `json:"file_id"`
	ID       string `json:"id"`
	Sz       int64  `json:"sz"`
	Resource []byte
	FromUser string `json:"from_user"`
}
