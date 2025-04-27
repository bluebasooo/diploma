package engine

// ID = userID
type History struct {
	ID     string
	UserID string
	Videos map[string]float64
}

func getHistory(id string) *History {
	return &History{
		ID:     id,
		Videos: make(map[string]float64),
	}
}

// mocked
func getHistories(ids []string) []History {
	histories := make([]History, 0, len(ids))
	for _, id := range ids {
		histories = append(histories, History{
			ID:     id,
			Videos: make(map[string]float64),
		})
	}
	return histories
}
