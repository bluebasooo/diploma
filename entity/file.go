package entity

type FilePart struct {
	ID       string `json:"id"`
	Sz       string `json:"sz"`
	Resource []byte
}
