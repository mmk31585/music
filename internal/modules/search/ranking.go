package search

import (
	"context"
	"time"
)

type RankingWeights struct {
	Popularity      float64
	Freshness       float64
	GenreMatch      float64
	Collaborative   float64
	Personalization float64
}

func DefaultWeights() RankingWeights {
	return RankingWeights{
		Popularity:      0.4,
		Freshness:       0.2,
		GenreMatch:      0.2,
		Collaborative:   0.15,
		Personalization: 0.05,
	}
}

type PersonalizationTerms struct {
	Genres    []string
	ArtistIDs []string
}

func buildRankedQuery(ctx context.Context, query, expandedQuery string, fields []string, userID string, weights RankingWeights) map[string]any {
	return buildRankedQueryWithPersonalization(ctx, query, expandedQuery, fields, userID, weights, nil)
}

func buildRankedQueryWithPersonalization(ctx context.Context, query, expandedQuery string, fields []string, userID string, weights RankingWeights, personalization *PersonalizationTerms) map[string]any {
	should := buildMultiMatchClauses(query, expandedQuery, fields)

	functions := []any{}

	functions = append(functions, map[string]any{
		"field_value_factor": map[string]any{
			"field":    "popularity",
			"factor":   weights.Popularity,
			"modifier": "log1p",
		},
	})

	functions = append(functions, map[string]any{
		"gauss": map[string]any{
			"created_at": map[string]any{
				"origin": time.Now().UTC().Format(time.RFC3339),
				"scale":  "90d",
				"decay":  0.5,
			},
		},
		"weight": weights.Freshness,
	})

	if personalization != nil && weights.Personalization > 0 {
		if len(personalization.Genres) > 0 {
			functions = append(functions, map[string]any{
				"filter": map[string]any{
					"terms": map[string]any{"genre": personalization.Genres},
				},
				"weight": weights.Personalization,
			})
		}
		if len(personalization.ArtistIDs) > 0 {
			functions = append(functions, map[string]any{
				"filter": map[string]any{
					"terms": map[string]any{"artist_id": personalization.ArtistIDs},
				},
				"weight": weights.Personalization,
			})
		}
	}

	return map[string]any{
		"size": 20,
		"query": map[string]any{
			"function_score": map[string]any{
				"query": map[string]any{
					"bool": map[string]any{
						"should":               should,
						"minimum_should_match": 1,
					},
				},
				"functions":  functions,
				"score_mode": "sum",
				"boost_mode": "multiply",
			},
		},
	}
}

func buildMultiMatchClauses(query, originalQuery string, fields []string) []any {
	should := []any{
		map[string]any{
			"multi_match": map[string]any{
				"query":     query,
				"fields":    fields,
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		},
	}

	for _, f := range fields {
		pf := "persian_" + f
		should = append(should, map[string]any{
			"multi_match": map[string]any{
				"query":  query,
				"fields": []string{pf + "^4"},
				"type":   "best_fields",
				"boost":  2.0,
			},
		})
	}

	acFields := []string{}
	for _, f := range fields {
		acFields = append(acFields, f+".autocomplete^3", f+".latin^2")
	}
	should = append(should, map[string]any{
		"multi_match": map[string]any{
			"query":  originalQuery,
			"fields": acFields,
			"type":   "bool_prefix",
		},
	})

	return should
}
