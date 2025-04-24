package engine

type Dot struct {
	ID       string
	History  map[string]float64 // history about 30 days
	BucketID string
}

type PlainHistory struct {
	DotID  string
	Videos map[string]DotVideoMetrics
}

type DotVideoMetrics struct {
	Views       int
	Likes       int
	Dislikes    int
	WatchTime   float64 // in minutes
	Engagement  float64 // engagement rate
	Comments    int
	ShareCount  int
	AverageView float64 // average view duration
}

type Bucket struct {
	ID                 string
	DotsToDistToCenter map[string]float64
	Center             map[string]float64
}

type VideoPool struct {
	BucketID string
	Videos   []string
}
