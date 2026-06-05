package config

type CORSConfig struct {
	AllowedOrigins []string
}

func loadCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: getEnvAsStringSlice("CORS_ALLOWED_ORIGINS", []string{
			"http://localhost:3000",
			"http://localhost:5173",
		}),
	}
}
