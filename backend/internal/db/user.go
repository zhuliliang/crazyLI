package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents an authenticated account.
type User struct {
	ID           uuid.UUID `json:"id"`
	GoogleID     string    `json:"google_id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Picture      string    `json:"picture"`
	PasswordHash string    `json:"-"`
	Provider     string    `json:"provider"`
	CreatedAt    time.Time `json:"created_at"`
}

// Store wraps common DB queries.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore creates a Store instance.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// UpsertUser ensures a Google user exists and returns the latest record.
func (s *Store) UpsertUser(ctx context.Context, u User) (User, error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now().UTC()

	const query = `
        INSERT INTO users (id, google_id, email, name, picture, provider, created_at)
        VALUES ($1, $2, $3, $4, $5, 'google', $6)
        ON CONFLICT (google_id)
        DO UPDATE SET email = EXCLUDED.email,
                      name = EXCLUDED.name,
                      picture = EXCLUDED.picture,
                      provider = 'google'
        RETURNING id, google_id, email, name, picture, password_hash, provider, created_at
    `

	return scanUser(s.pool.QueryRow(ctx, query, u.ID, u.GoogleID, u.Email, u.Name, u.Picture, u.CreatedAt))
}

// CreateEmailUser inserts a new email/password user.
func (s *Store) CreateEmailUser(ctx context.Context, u User) (User, error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now().UTC()

	const query = `
        INSERT INTO users (id, email, name, password_hash, provider, created_at)
        VALUES ($1, $2, $3, $4, 'email', $5)
        RETURNING id, google_id, email, name, picture, password_hash, provider, created_at
    `

	return scanUser(s.pool.QueryRow(ctx, query, u.ID, u.Email, u.Name, u.PasswordHash, u.CreatedAt))
}

// GetUserByEmail fetches by email.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (User, error) {
	const query = `SELECT id, google_id, email, name, picture, password_hash, provider, created_at FROM users WHERE email = $1`
	return scanUser(s.pool.QueryRow(ctx, query, email))
}

// GetUserByID fetches a user by UUID.
func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	const query = `SELECT id, google_id, email, name, picture, password_hash, provider, created_at FROM users WHERE id = $1`
	return scanUser(s.pool.QueryRow(ctx, query, id))
}

func scanUser(row pgx.Row) (User, error) {
	var result User
	var googleID sql.NullString
	var picture sql.NullString
	var passwordHash sql.NullString
	if err := row.Scan(&result.ID, &googleID, &result.Email, &result.Name, &picture, &passwordHash, &result.Provider, &result.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	result.GoogleID = googleID.String
	result.Picture = picture.String
	result.PasswordHash = passwordHash.String
	return result, nil
}
