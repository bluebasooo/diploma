package entity

// storage
type FilePart struct {
	ID       string `json:"id"`
	Sz       int64  `json:"sz"`
	Resource []byte
}
