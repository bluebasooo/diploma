package models

type WritePlan struct {
	Ops []Operations
}

type Operations struct {
	HashOperation string
	BytesFrom     int64
	BytesTo       int64
}
