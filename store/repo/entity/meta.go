package entity

import "time"

// mongo
type FileMeta struct {
	ID           string         `bson:"_id"` // set by mongodb
	Parts        []FileMetaPart `bson:"parts"`
	FullSz       int64          `bson:"fullSz"`
	CreatedAt    time.Time      `bson:"createdAt"`
	UpdatedAt    time.Time      `bson:"updatedAt"`
	PartSequence []string       `bson:"partSequence"`
	IsDraft      bool           `bson:"isDraft"`
}

type FileMetaPart struct {
	Hash  string `bson:"hash"`
	Sz    int64  `bson:"sz"`
	S3Url string `bson:"s3Url"`
}
