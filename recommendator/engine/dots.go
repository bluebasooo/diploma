package engine

import "dev/bluebasooo/video-recomendator/service/model"

type Dot struct {
	ID       string
	History  map[string]float64 // history about 30 days
	BucketID string
}

func getDot(id string) *Dot { // tmp
	return &Dot{
		ID:       id,
		History:  make(map[string]float64),
		BucketID: "",
	}
}

func getDots(ids []string) []Dot {
	dots := make([]Dot, 0, len(ids))
	for _, id := range ids {
		dots = append(dots, *getDot(id))
	}
	return dots
}

func fromHistory(history []model.History) *Dot 
	
	
	return dot
}

type PlainHistory struct {
	DotID  string
	Videos map[string]DotVideoMetrics
}

type VideoPool struct {
	BucketID string
	Videos   []string
}
