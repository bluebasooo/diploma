package search

import (
	"bytes"
	"encoding/json"
	"log"
)

type Body = interface{}

type BulkRequest struct {
	Instructions []Instruction
}

func (r *BulkRequest) toReader() *bytes.Reader {
	var buf bytes.Buffer
	for _, doc := range r.Instructions {
		meta, err := json.Marshal(doc.Meta)
		if err != nil {
			log.Fatalf("Error marshaling document %s: %s", doc.Meta.Index.ID, err)
		}
		meta = append(meta, byte('\n'))

		data, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Error marshaling document %s: %s", doc.Meta.Index.ID, err)
		}
		data = append(data, byte('\n'))

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	return bytes.NewReader(buf.Bytes())
}

type Instruction struct {
	Meta Meta
	Body Body
}

type Meta struct {
	Index  *IndexPreamble    `json:"index"`
	Delete *IndexPreamble    `json:"delete"`
	Create *IndexPreamble    `json:"create"`
	Update *DocIndexPreamble `json:"update"`
}

func (m *Meta) IndexName() string {
	if m.Index != nil {
		return m.Index.Index
	}
	if m.Create != nil {
		return m.Create.Index
	}
	if m.Update != nil {
		return m.Update.Doc.Index
	}
	if m.Delete != nil {
		return m.Delete.Index
	}
	return ""
}

type DocIndexPreamble struct {
	Doc IndexPreamble `json:"doc"`
}

type IndexPreamble struct {
	Index string `json:"_index"`
	ID    string `json:"_id"`
}

type SearchResponse[T any] struct {
	took      int  `json:"took"`
	timed_out bool `json:"timed_out"`
	Shards    struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source T       `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
