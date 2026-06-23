//go:build pgx_integration

package auth

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	now := time.Now()
	userID := uuid.New().String()
	email := "test@example.com"
	username := "testuser"
	displayName := "Test User"
	passwordHash := "$2a$10$hashedpassword"

	rows := sqlmock.NewRows([]string{
		"id", "email", "username", "display_name", "password_hash",
		"avatar_url", "role", "is_active", "email_verified_at", "created_at", "updated_at",
	}).AddRow(
		userID, email, username, displayName, passwordHash,
		nil, "listener", true, nil, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (email, username, display_name, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at
	`)).WithArgs(email, username, displayName, passwordHash).
		WillReturnRows(rows)

	user, err := repo.CreateUser(context.Background(), User{
		Email:        email,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
	})
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "listener", user.Role)
	assert.True(t, user.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	userID := uuid.New().String()
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{
		"id", "email", "username", "display_name", "password_hash",
		"avatar_url", "role", "is_active", "email_verified_at", "created_at", "updated_at",
	}).AddRow(
		userID, email, "testuser", "Test User", "$2a$10$hash",
		nil, "listener", true, nil, time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL`)).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "testuser", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, username, display_name, password_hash, avatar_url, role, is_active, email_verified_at, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL`)).
		WithArgs("nonexistent@example.com").
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetUserByEmail(context.Background(), "nonexistent@example.com")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	userID := uuid.New().String()
	tokenHash := "sha256hash"
	familyID := uuid.New().String()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO refresh_tokens (user_id, token_hash, family_id, expires_at)
		VALUES ($1, $2, $3, $4)
	`)).WithArgs(userID, tokenHash, familyID, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateRefreshToken(context.Background(), userID, tokenHash, familyID, expiresAt)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
