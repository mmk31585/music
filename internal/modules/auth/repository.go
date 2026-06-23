package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	apperrors "music/internal/common/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *Repository) CreateUser(ctx context.Context, user User) (User, error) {
	query := `
		INSERT INTO users (email, username, display_name, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at
	`

	var created User

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.Username,
		user.DisplayName,
		user.PasswordHash,
	).Scan(
		&created.ID,
		&created.Email,
		&created.Username,
		&created.DisplayName,
		&created.PasswordHash,
		&created.AvatarURL,
		&created.Role,
		&created.IsActive,
		&created.EmailVerifiedAt,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		if isUniqueViolation(err) {
			return User{}, apperrors.Conflict("email or username already exists", nil)
		}

		return User{}, err
	}

	return created, nil
}

func (r *Repository) FindUserByEmailOrUsername(ctx context.Context, value string) (User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE email = $1 OR username = $1
		LIMIT 1
	`

	return r.scanUser(r.db.QueryRow(ctx, query, value))
}

func (r *Repository) FindUserByID(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	return r.scanUser(r.db.QueryRow(ctx, query, id))
}

func (r *Repository) CreateSession(
	ctx context.Context,
	userID string,
	refreshToken string,
	userAgent *string,
	ipAddress *string,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO auth_sessions (user_id, refresh_token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		userID,
		HashToken(refreshToken),
		userAgent,
		ipAddress,
		expiresAt,
	)

	return err
}

func (r *Repository) FindValidSessionByRefreshToken(ctx context.Context, refreshToken string) (AuthSession, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, user_agent, ip_address, expires_at, revoked_at, created_at
		FROM auth_sessions
		WHERE refresh_token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
		LIMIT 1
	`

	var session AuthSession

	err := r.db.QueryRow(ctx, query, HashToken(refreshToken)).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return AuthSession{}, apperrors.Unauthorized("invalid refresh token", nil)
	}

	if err != nil {
		return AuthSession{}, err
	}

	return session, nil
}

func (r *Repository) RevokeSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = NOW()
		WHERE refresh_token_hash = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, HashToken(refreshToken))
	return err
}

func (r *Repository) RevokeSessionByID(ctx context.Context, sessionID string) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = NOW()
		WHERE id = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, sessionID)
	return err
}

func (r *Repository) scanUser(row pgx.Row) (User, error) {
	var user User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.DisplayName,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.Role,
		&user.IsActive,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, apperrors.NotFound("user not found", nil)
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}
func (r *Repository) FindPublicUser(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
		LIMIT 1
	`

	return r.scanUser(r.db.QueryRow(ctx, query, id))
}

func (r *Repository) ListUsers(ctx context.Context, params ListUsersParams) ([]User, int, error) {
	where := "WHERE 1=1"
	args := make([]any, 0)
	argIdx := 1

	if params.Search != "" {
		where += fmt.Sprintf(" AND (email ILIKE $%d OR username ILIKE $%d OR display_name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	var countQuery = "SELECT COUNT(*) FROM users " + where
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "created_at DESC"
	sortByMap := map[string]string{
		"display_name": "display_name",
		"email":        "email",
		"username":     "username",
		"role":         "role",
		"is_active":    "is_active",
		"created_at":   "created_at",
	}
	if col, ok := sortByMap[params.SortBy]; ok {
		order := "DESC"
		if params.SortOrder == "asc" || params.SortOrder == "ASC" {
			order = "ASC"
		}
		orderBy = col + " " + order
	}

	offset := (params.Page - 1) * params.PageSize
	limit := params.PageSize

	query := fmt.Sprintf(`
		SELECT id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at
		FROM users %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, where, orderBy, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.DisplayName,
			&user.PasswordHash, &user.AvatarURL, &user.Role, &user.IsActive,
			&user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// allowedUserUpdateColumns is the strict whitelist of column names that can be
// dynamically referenced in UPDATE SET clauses. Any key not in this map is
// rejected to prevent SQL injection via uncontrolled column names.
var allowedUserUpdateColumns = map[string]string{
	"role":           "role",
	"is_active":      "is_active",
	"email_verified": "email_verified_at",
	"display_name":   "display_name",
	"username":       "username",
	"email":          "email",
	"avatar_url":     "avatar_url",
	"bio":            "bio",
	"location":       "location",
	"website":        "website",
	"preferences":    "preferences",
}

func (r *Repository) UpdateUser(ctx context.Context, id string, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]any, 0, len(updates)+1)
	argIdx := 1

	for key, value := range updates {
		col, ok := allowedUserUpdateColumns[key]
		if !ok {
			// Reject unknown keys instead of interpolating them raw (SQL injection prevention)
			return apperrors.BadRequest("unknown user field: "+key, map[string]string{
				"field":   key,
				"allowed": "role, is_active, email_verified, display_name, username, email, avatar_url",
			})
		}

		// Handle email_verified as a special boolean→SET NULL logic
		if key == "email_verified" {
			if b, ok := value.(bool); ok && b {
				setClauses = append(setClauses, fmt.Sprintf("email_verified_at = NOW()"))
				continue
			}
			// Set to NULL when false
			setClauses = append(setClauses, fmt.Sprintf("email_verified_at = $%d", argIdx))
			args = append(args, nil)
			argIdx++
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, value)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, passwordHash, id)
	return err
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
