package dto

type FileMetaPlanDto struct {
	SizeInBytes int64 `json:"sizeInBytes"`
}

type FileMetaDto struct {
	ID     string            `json:"id"`
	Parts  []FileMetaPartDto `json:"parts"`
	FullSz int64             `json:"fullSz"`
}

type FileMetaPartDto struct {
	Hash  string `json:"hash"`
	Sz    int64  `json:"sz"`
	S3Url string `json:"s3Url"`
}
