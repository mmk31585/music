package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AuditEntry struct {
	ID           string          `db:"id"`
	UserID       *string         `db:"user_id"`
	Action       string          `db:"action"`
	ResourceType string          `db:"resource_type"`
	ResourceID   *string         `db:"resource_id"`
	RequestBody  json.RawMessage `db:"request_body"`
	IPAddress    string          `db:"ip_address"`
	UserAgent    string          `db:"user_agent"`
	CreatedAt    time.Time       `db:"created_at"`
}

type AuditLogger struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewAuditLogger(db *sqlx.DB, logger *zap.Logger) *AuditLogger {
	return &AuditLogger{db: db, logger: logger}
}

func (al *AuditLogger) Middleware(skipPaths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool, len(skipPaths))
	for _, p := range skipPaths {
		skipMap[p] = true
	}

	return func(c *gin.Context) {
		if skipMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		writer := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		userID, _ := c.Get("auth_user_id")
		var uid *string
		if u, ok := userID.(string); ok {
			uid = &u
		}

		entry := AuditEntry{
			ID:           uuid.New().String(),
			UserID:       uid,
			Action:       c.Request.Method + " " + c.Request.URL.Path,
			ResourceType: extractResourceType(c.Request.URL.Path),
			RequestBody:  bodyBytes,
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			CreatedAt:    time.Now().UTC(),
		}

		if len(c.Params) > 0 {
			for _, param := range c.Params {
				if param.Key == "id" || param.Key == "userId" || param.Key == "trackId" {
					v := param.Value
					entry.ResourceID = &v
					break
				}
			}
		}

		if err := al.insertAuditLog(c, entry); err != nil {
			al.logger.Warn("failed to write audit log", zap.Error(err))
		}
	}
}

func (al *AuditLogger) insertAuditLog(c *gin.Context, entry AuditEntry) error {
	_, err := al.db.ExecContext(c.Request.Context(),
		`INSERT INTO audit_log (id, user_id, action, resource_type, request_body, ip_address, user_agent, created_at)
		 VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7, $8)`,
		entry.ID, entry.UserID, entry.Action, entry.ResourceType,
		entry.RequestBody, entry.IPAddress, entry.UserAgent, entry.CreatedAt,
	)
	return err
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func extractResourceType(path string) string {
	parts := splitPath(path)
	for _, part := range parts {
		switch part {
		case "tracks", "albums", "artists", "playlists", "users",
			"auth", "media", "social", "subscription", "tips":
			return part
		}
	}
	return "unknown"
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, c := range path {
		if c == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
