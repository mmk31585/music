package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

const tracksMapping = `{
	"settings": {
		"number_of_shards": 3,
		"number_of_replicas": 1,
		"analysis": {
			"filter": {
				"persian_stop": {
					"type": "stop",
					"stopwords": ["_persian_"]
				},
				"persian_synonym": {
					"type": "synonym",
					"synonyms": [
						"pop, پاپ", "rock, راک", "jazz, جاز",
						"classical, کلاسیک", "traditional, سنتی",
						"folk, فولک, محلی", "electronic, الکترونیک",
						"hip-hop, هیپ_هاپ", "rap, رپ",
						"persian, فارسی, پارسی, ایرانی, irani"
					]
				},
				"persian_stemmer": {
					"type": "stemmer",
					"language": "arabic"
				}
			},
			"analyzer": {
				"persian_analyzer": {
					"type": "custom",
					"tokenizer": "standard",
					"filter": [
						"lowercase",
						"arabic_normalization",
						"persian_normalization",
						"persian_stop",
						"persian_synonym",
						"persian_stemmer"
					]
				},
				"persian_autocomplete": {
					"type": "custom",
					"tokenizer": "persian_edge_ngram",
					"filter": ["lowercase", "arabic_normalization", "persian_normalization"]
				},
				"latin_autocomplete": {
					"type": "custom",
					"tokenizer": "latin_edge_ngram",
					"filter": ["lowercase", "asciifolding"]
				}
			},
			"tokenizer": {
				"persian_edge_ngram": {
					"type": "edge_ngram",
					"min_gram": 2,
					"max_gram": 20,
					"token_chars": ["letter", "digit"]
				},
				"latin_edge_ngram": {
					"type": "edge_ngram",
					"min_gram": 2,
					"max_gram": 20,
					"token_chars": ["letter", "digit"]
				}
			}
		}
	},
	"mappings": {
		"dynamic": false,
		"properties": {
			"id": { "type": "keyword" },
			"title": {
				"type": "search_as_you_type",
				"analyzer": "persian_analyzer",
				"fields": {
					"autocomplete": { "type": "text", "analyzer": "persian_autocomplete" },
					"latin": { "type": "text", "analyzer": "latin_autocomplete" },
					"completion": { "type": "completion" }
				}
			},
			"persian_title": {
				"type": "text",
				"analyzer": "persian_analyzer",
				"fields": {
					"autocomplete": { "type": "text", "analyzer": "persian_autocomplete" },
					"completion": { "type": "completion" }
				}
			},
			"artist_name": {
				"type": "text",
				"analyzer": "persian_analyzer",
				"fields": {
					"autocomplete": { "type": "text", "analyzer": "persian_autocomplete" },
					"latin": { "type": "text", "analyzer": "latin_autocomplete" },
					"completion": { "type": "completion" }
				}
			},
			"album_title": {
				"type": "text",
				"analyzer": "persian_analyzer",
				"fields": {
					"autocomplete": { "type": "text", "analyzer": "persian_autocomplete" }
				}
			},
			"lyrics": { "type": "text", "analyzer": "persian_analyzer" },
			"genre_names": { "type": "keyword" },
			"duration": { "type": "integer" },
			"year": { "type": "integer" },
			"popularity": { "type": "float" },
			"is_explicit": { "type": "boolean" },
			"created_at": { "type": "date" },
			"similar_tracks": {
				"type": "nested",
				"properties": {
					"track_id": { "type": "keyword" },
					"similarity": { "type": "float" }
				}
			}
		}
	}
}`

const artistsMapping = `{
	"settings": {
		"number_of_shards": 2,
		"number_of_replicas": 1,
		"analysis": {
			"filter": {
				"persian_stop": {
					"type": "stop",
					"stopwords": ["_persian_"]
				},
				"persian_synonym": {
					"type": "synonym",
					"synonyms": [
						"pop, پاپ", "rock, راک", "jazz, جاز",
						"classical, کلاسیک", "traditional, سنتی",
						"persian, فارسی, پارسی, ایرانی, irani"
					]
				}
			},
			"analyzer": {
				"persian_analyzer": {
					"type": "custom",
					"tokenizer": "standard",
					"filter": ["lowercase", "arabic_normalization", "persian_normalization", "persian_stop", "persian_synonym"]
				},
				"persian_autocomplete": {
					"type": "custom",
					"tokenizer": "persian_edge_ngram",
					"filter": ["lowercase", "arabic_normalization", "persian_normalization"]
				},
				"latin_autocomplete": {
					"type": "custom",
					"tokenizer": "latin_edge_ngram",
					"filter": ["lowercase", "asciifolding"]
				}
			},
			"tokenizer": {
				"persian_edge_ngram": { "type": "edge_ngram", "min_gram": 2, "max_gram": 20, "token_chars": ["letter", "digit"] },
				"latin_edge_ngram": { "type": "edge_ngram", "min_gram": 2, "max_gram": 20, "token_chars": ["letter", "digit"] }
			}
		}
	},
	"mappings": {
		"dynamic": false,
		"properties": {
			"id": { "type": "keyword" },
			"name": {
				"type": "search_as_you_type",
				"analyzer": "persian_analyzer",
				"fields": {
					"autocomplete": { "type": "text", "analyzer": "persian_autocomplete" },
					"latin": { "type": "text", "analyzer": "latin_autocomplete" },
					"completion": { "type": "completion" }
				}
			},
			"persian_name": {
				"type": "text",
				"analyzer": "persian_analyzer",
				"fields": {
					"autocomplete": { "type": "text", "analyzer": "persian_autocomplete" }
				}
			},
			"bio": { "type": "text", "analyzer": "persian_analyzer" },
			"persian_bio": { "type": "text", "analyzer": "persian_analyzer" },
			"monthly_listeners": { "type": "integer" },
			"is_verified": { "type": "boolean" },
			"genre_names": { "type": "keyword" }
		}
	}
}`

const albumsMapping = `{
	"settings": {
		"number_of_shards": 2,
		"number_of_replicas": 1,
		"analysis": {
			"filter": {
				"persian_stop": { "type": "stop", "stopwords": ["_persian_"] }
			},
			"analyzer": {
				"persian_analyzer": {
					"type": "custom",
					"tokenizer": "standard",
					"filter": ["lowercase", "arabic_normalization", "persian_normalization", "persian_stop"]
				}
			}
		}
	},
	"mappings": {
		"dynamic": false,
		"properties": {
			"id": { "type": "keyword" },
			"title": { "type": "text", "analyzer": "persian_analyzer" },
			"persian_title": { "type": "text", "analyzer": "persian_analyzer" },
			"artist_name": { "type": "text", "analyzer": "persian_analyzer" },
			"release_year": { "type": "integer" },
			"album_type": { "type": "keyword" },
			"genre_names": { "type": "keyword" }
		}
	}
}`

type IndexDef struct {
	Name    string
	Mapping string
	Alias   string
}

func IndexDefinitions() []IndexDef {
	return []IndexDef{
		{Name: "tracks-v1", Mapping: tracksMapping, Alias: "tracks"},
		{Name: "artists-v1", Mapping: artistsMapping, Alias: "artists"},
		{Name: "albums-v1", Mapping: albumsMapping, Alias: "albums"},
	}
}

func EnsureIndexes(ctx context.Context, client *Client, logger *zap.Logger) error {
	for _, idx := range IndexDefinitions() {
		exists, err := client.IndexExists(ctx, idx.Alias)
		if err != nil {
			return fmt.Errorf("check index %s: %w", idx.Alias, err)
		}
		if exists {
			logger.Debug("index exists", zap.String("index", idx.Alias))
			continue
		}

		if err := client.CreateIndex(ctx, idx.Name, idx.Mapping); err != nil {
			return fmt.Errorf("create index %s: %w", idx.Name, err)
		}

		if err := putAlias(ctx, client, idx.Name, idx.Alias); err != nil {
			return fmt.Errorf("alias %s -> %s: %w", idx.Alias, idx.Name, err)
		}

		logger.Info("created opensearch index with alias",
			zap.String("index", idx.Name),
			zap.String("alias", idx.Alias),
		)
	}
	return nil
}

func putAlias(ctx context.Context, client *Client, index, alias string) error {
	body := fmt.Sprintf(`{"actions":[{"add":{"index":"%s","alias":"%s"}}]}`, index, alias)

	resp, err := client.Client().Indices.UpdateAliases(
		strings.NewReader(body),
		client.Client().Indices.UpdateAliases.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("alias error: %s", string(b))
	}

	return nil
}

func (c *Client) UpdateIndexMapping(ctx context.Context, index string, mapping map[string]any) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	req := opensearchapi.IndicesPutMappingRequest{
		Index: []string{index},
		Body:  strings.NewReader(string(body)),
	}

	resp, err := req.Do(ctx, c.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update mapping error: %s", string(b))
	}

	return nil
}
