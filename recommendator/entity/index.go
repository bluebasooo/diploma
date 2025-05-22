package entity

import "time"

// metadata about last processed data
type Index struct {
	ID          int64     `json:"id"`
	LastUpdated time.Time `json:"last_updated"`
}
