package config

import "os"

type OpenSearchConfig struct {
	URL      string
	Username string
	Password string
	Index    string
}

func loadOpenSearchConfig() OpenSearchConfig {
	index := os.Getenv("OPENSEARCH_INDEX")
	if index == "" {
		index = "tracks"
	}
	return OpenSearchConfig{
		URL:      os.Getenv("OPENSEARCH_URL"),
		Username: os.Getenv("OPENSEARCH_USERNAME"),
		Password: os.Getenv("OPENSEARCH_PASSWORD"),
		Index:    index,
	}
}
