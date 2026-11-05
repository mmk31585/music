package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zaptest"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		allowedOrigins     []string
		requestOrigin      string
		requestMethod      string
		requestHeaders     string
		expectedStatus     int
		expectACAO         string
		expectCredentials  string
		expectVary         string
		expectAllowHeaders string
		expectMaxAge       string
	}{
		{
			name:               "allowed origin preflight",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "https://example.com",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "Authorization,Content-Type",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "https://example.com",
			expectCredentials:  "true",
			expectVary:         "Origin, Access-Control-Request-Headers",
			expectAllowHeaders: "Authorization,Content-Type",
			expectMaxAge:       "600",
		},
		{
			name:               "unlisted origin preflight - no ACAO",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "https://evil.com",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "Authorization,Content-Type",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "",
			expectCredentials:  "",
			expectVary:         "",
			expectAllowHeaders: "",
			expectMaxAge:       "",
		},
		{
			name:               "trailing slash in config",
			allowedOrigins:     []string{"https://example.com/"},
			requestOrigin:      "https://example.com",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "Authorization",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "https://example.com",
			expectCredentials:  "true",
			expectVary:         "Origin, Access-Control-Request-Headers",
			expectAllowHeaders: "Authorization",
			expectMaxAge:       "600",
		},
		{
			name:               "Authorization in requested headers",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "https://example.com",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "Authorization,X-Custom-Header",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "https://example.com",
			expectCredentials:  "true",
			expectVary:         "Origin, Access-Control-Request-Headers",
			expectAllowHeaders: "Authorization,X-Custom-Header",
			expectMaxAge:       "600",
		},
		{
			name:               "normal request allowed origin carries CORS headers",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "https://example.com",
			requestMethod:      http.MethodGet,
			requestHeaders:     "",
			expectedStatus:     http.StatusOK,
			expectACAO:         "https://example.com",
			expectCredentials:  "true",
			expectVary:         "Origin",
			expectAllowHeaders: "",
			expectMaxAge:       "",
		},
		{
			name:               "500 response still carries CORS headers",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "https://example.com",
			requestMethod:      http.MethodGet,
			requestHeaders:     "",
			expectedStatus:     http.StatusInternalServerError,
			expectACAO:         "https://example.com",
			expectCredentials:  "true",
			expectVary:         "Origin",
			expectAllowHeaders: "",
			expectMaxAge:       "",
		},
		{
			name:               "fallback headers when no Access-Control-Request-Headers",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "https://example.com",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "https://example.com",
			expectCredentials:  "true",
			expectVary:         "Origin, Access-Control-Request-Headers",
			expectAllowHeaders: "Authorization,Content-Type,Accept,X-Request-ID,X-Vercel-Oidc-Token",
			expectMaxAge:       "600",
		},
		{
			name:               "multiple allowed origins",
			allowedOrigins:     []string{"https://example.com", "https://app.example.com"},
			requestOrigin:      "https://app.example.com",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "Authorization",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "https://app.example.com",
			expectCredentials:  "true",
			expectVary:         "Origin, Access-Control-Request-Headers",
			expectAllowHeaders: "Authorization",
			expectMaxAge:       "600",
		},
		{
			name:               "no origin header on preflight - passes through",
			allowedOrigins:     []string{"https://example.com"},
			requestOrigin:      "",
			requestMethod:      http.MethodOptions,
			requestHeaders:     "Authorization",
			expectedStatus:     http.StatusNoContent,
			expectACAO:         "",
			expectCredentials:  "",
			expectVary:         "",
			expectAllowHeaders: "",
			expectMaxAge:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)

			// Build router
			router := gin.New()
			router.Use(CORS(tt.allowedOrigins, logger))

			if tt.expectedStatus == http.StatusInternalServerError {
				router.GET("/test", func(c *gin.Context) {
					c.AbortWithStatus(http.StatusInternalServerError)
				})
			} else {
				router.GET("/test", func(c *gin.Context) {
					c.Status(http.StatusOK)
				})
			}
			router.OPTIONS("/test", func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})

			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}
			if tt.requestHeaders != "" {
				req.Header.Set("Access-Control-Request-Headers", tt.requestHeaders)
			}
			if tt.requestMethod == http.MethodOptions {
				req.Header.Set("Access-Control-Request-Method", http.MethodGet)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectACAO != "" {
				if w.Header().Get("Access-Control-Allow-Origin") != tt.expectACAO {
					t.Errorf("expected ACAO %q, got %q", tt.expectACAO, w.Header().Get("Access-Control-Allow-Origin"))
				}
			} else {
				if w.Header().Get("Access-Control-Allow-Origin") != "" {
					t.Errorf("expected no ACAO, got %q", w.Header().Get("Access-Control-Allow-Origin"))
				}
			}

			if tt.expectCredentials != "" {
				if w.Header().Get("Access-Control-Allow-Credentials") != tt.expectCredentials {
					t.Errorf("expected credentials %q, got %q", tt.expectCredentials, w.Header().Get("Access-Control-Allow-Credentials"))
				}
			} else {
				if w.Header().Get("Access-Control-Allow-Credentials") != "" {
					t.Errorf("expected no credentials, got %q", w.Header().Get("Access-Control-Allow-Credentials"))
				}
			}

			if tt.expectVary != "" {
				if !strings.Contains(w.Header().Get("Vary"), tt.expectVary) {
					t.Errorf("expected Vary to contain %q, got %q", tt.expectVary, w.Header().Get("Vary"))
				}
			} else {
				if w.Header().Get("Vary") != "" {
					t.Errorf("expected no Vary, got %q", w.Header().Get("Vary"))
				}
			}

			if tt.expectAllowHeaders != "" {
				if w.Header().Get("Access-Control-Allow-Headers") != tt.expectAllowHeaders {
					t.Errorf("expected Allow-Headers %q, got %q", tt.expectAllowHeaders, w.Header().Get("Access-Control-Allow-Headers"))
				}
			} else {
				if tt.requestMethod == http.MethodOptions {
					if w.Header().Get("Access-Control-Allow-Headers") != "" {
						t.Errorf("expected no Allow-Headers, got %q", w.Header().Get("Access-Control-Allow-Headers"))
					}
				}
			}

			if tt.expectMaxAge != "" {
				if w.Header().Get("Access-Control-Max-Age") != tt.expectMaxAge {
					t.Errorf("expected Max-Age %q, got %q", tt.expectMaxAge, w.Header().Get("Access-Control-Max-Age"))
				}
			} else {
				if tt.requestMethod == http.MethodOptions {
					if w.Header().Get("Access-Control-Max-Age") != "" {
						t.Errorf("expected no Max-Age, got %q", w.Header().Get("Access-Control-Max-Age"))
					}
				}
			}

			// Preflight should not call Next() when origin is not allowed
			if tt.requestOrigin != "" && tt.requestMethod == http.MethodOptions {
				matched := false
				for _, o := range tt.allowedOrigins {
					o = strings.TrimSpace(o)
					o = strings.TrimRight(o, "/")
					if o == tt.requestOrigin {
						matched = true
						break
					}
				}
				if !matched && w.Code != http.StatusNoContent {
					t.Errorf("unlisted origin should get 204 without CORS headers")
				}
			}
		})
	}
}

func TestCORS_PanicRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := zaptest.NewLogger(t)

	router := gin.New()
	router.Use(CORS([]string{"https://example.com"}, logger))
	router.Use(gin.Recovery())
	router.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("Origin", "https://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should recover and return 500 with CORS headers
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("expected ACAO on panic response, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Errorf("expected credentials on panic response")
	}
}