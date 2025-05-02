package repo

import (
	"dev/bluebasooo/video-platform/repo/entity"
	"dev/bluebasooo/video-platform/search"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type SearchRepo struct {
	elasticDB *search.ElasticDB
}

func NewSearchRepo(elasticDB *search.ElasticDB) *SearchRepo {
	repo := &SearchRepo{elasticDB: elasticDB}
	repo.init()
	return repo
}

func (r *SearchRepo) init() {
	r.initIndex(videoIndexName, videoMapping)
	r.initIndex(authorIndexName, authorMapping)
}

func (r *SearchRepo) initIndex(indexName string, mapping string) {
	es := r.elasticDB.GetClient()

	resp, err := es.Indices.Create(indexName, es.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		log.Fatalf("Error creating index: %s", resp.Status())
	}
}

func (r *SearchRepo) IndexVideo(videos map[string]entity.VideoIndex) error {
	entities := make([]search.Instruction, 0, len(videos))

	for id, video := range videos {
		meta := search.Meta{
			Index: &search.IndexPreamble{
				Index: videoIndexName,
				ID:    id,
			},
		}

		entities = append(entities, search.Instruction{
			Meta: meta,
			Body: video,
		})
	}

	return r.elasticDB.BulkIndexEntities(entities)
}

func (r *SearchRepo) IndexAuthor(authors map[string]entity.AuthorIndex) error {
	entities := make([]search.Instruction, 0, len(authors))

	for id, author := range authors {
		meta := search.Meta{
			Index: &search.IndexPreamble{
				Index: authorIndexName,
				ID:    id,
			},
		}

		entities = append(entities, search.Instruction{
			Meta: meta,
			Body: author,
		})
	}

	return r.BulkIndexEntities(entities)
}

func (r *SearchRepo) BulkIndexEntities(entities []search.Instruction) error {
	return r.elasticDB.BulkIndexEntities(entities)
}

func (r *SearchRepo) SearchVideos(searched string) ([]entity.VideoSearchResult, error) {
	es := r.elasticDB.GetClient()

	query := fmt.Sprintf(defaultSearchVideoQuery, searched, searched, searched)

	resp, err := es.Search(
		es.Search.WithIndex(videoIndexName),
		es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("error searching videos: %s", resp.Status())
	}

	var decodedResponse search.SearchResponse[entity.VideoSearchResult]
	if err := json.NewDecoder(resp.Body).Decode(&decodedResponse); err != nil {
		return nil, err
	}

	var result []entity.VideoSearchResult
	for _, hit := range decodedResponse.Hits.Hits {
		result = append(result, hit.Source)
	}

	return result, nil
}

func (r *SearchRepo) SearchAuthors(query string) ([]entity.AuthorSearchResult, error) {
	es := r.elasticDB.GetClient()

	searchQuery := fmt.Sprintf(defaultSearchAuthorQuery, query, query, query)

	resp, err := es.Search(
		es.Search.WithIndex(authorIndexName),
		es.Search.WithBody(strings.NewReader(searchQuery)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("error searching authors: %s", resp.Status())
	}

	var decodedResponse search.SearchResponse[entity.AuthorSearchResult]
	if err := json.NewDecoder(resp.Body).Decode(&decodedResponse); err != nil {
		return nil, err
	}

	var result []entity.AuthorSearchResult
	for _, hit := range decodedResponse.Hits.Hits {
		result = append(result, hit.Source)
	}
	return result, nil
}

var defaultSearchAuthorQuery = `
{
	"query": {
		"bool": {
			"should": [
				{
					"multi_match": {
						"query": "%s",
						"fields": ["name^2", "username^4"],
						"type": "best_fields",
						"operator": "or"
					}
				},
				{
					"match_phrase": {
						"username": "%s",
						"boost": 4
					}
				},
				{
					"match_phrase": {
						"name": "%s",
						"boost": 2
					}
				}
			],
			"minimum_should_match": 1
		}
	}
}
`

var defaultSearchVideoQuery = `
{
  "query": {
    "function_score": {
      "query": {
        "bool": {
          "should": [
            {
              "multi_match": {
                "query": "%s",
                "fields": ["title^4", "authorName^3", "description^2"],
                "type": "best_fields",
                "operator": "or"
              }
            },
            {
              "match_phrase": {
                "title": {
                  "query": "%s",
                  "boost": 4
                }
              }
            },
            {
              "match_phrase": {
                "authorName": {
                  "query": "%s",
                  "boost": 3
                }
              }
            }
          ],
          "minimum_should_match": 1
        }
      },
      "boost_mode": "multiply",
      "score_mode": "sum",
      "functions": [
        {
          "field_value_factor": {
            "field": "views",
            "factor": 0.0001,
            "modifier": "sqrt",
            "missing": 1
          }
        },
        {
          "gauss": {
            "uploadDate": {
              "origin": "now",
              "scale": "30d",
              "decay": 0.5
            }
          }
        }
      ]
    }
  }
}
`

var videoIndexName = "video_index"
var authorIndexName = "author_index"

var videoMapping = `
{
	"mappings": {
		"properties": {
			"title": {
				"type": "text"
			},
			"description": {
				"type": "text"
			},
			"durationMs": {
				"type": "integer"
			},
			"authorName": {
				"type": "text"
			},
			"uploadDate": {
				"type": "date"
			},
			"views": {
				"type": "integer"
			},
			"hidden": {
				"type": "boolean"
			}
		}
	}
}
`

var authorMapping = `
{
	"mappings": {
		"properties": {
			"name": {
				"type": "text"
			},
			"username": {
				"type": "text"
			}
		}
	}
}
`
