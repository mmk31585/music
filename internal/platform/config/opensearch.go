package config

type OpenSearchConfig struct {
	URL      string
	Username string
	Password string
	Index    string
}

func loadOpenSearchConfig() OpenSearchConfig {
	return OpenSearchConfig{
		URL:      getEnv("OPENSEARCH_URL", ""),
		Username: getEnv("OPENSEARCH_USERNAME", ""),
		Password: getEnv("OPENSEARCH_PASSWORD", ""),
		Index:    getEnv("OPENSEARCH_INDEX", "tracks"),
	}
}
