package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	apperrors "music/internal/contracts/errors"

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
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
