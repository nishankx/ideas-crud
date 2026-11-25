package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nisha/gravehard/internal/domain"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash, displayName string, isDemo bool) (*domain.User, error) {
	const q = `
		INSERT INTO users (email, password_hash, display_name, is_demo)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password_hash, display_name, is_demo, created_at`

	var u domain.User
	err := r.db.QueryRow(ctx, q, email, passwordHash, displayName, isDemo).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.IsDemo, &u.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrConflict
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, display_name, is_demo, created_at
		FROM users
		WHERE lower(email) = lower($1)
		LIMIT 1`

	var u domain.User
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.IsDemo, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, display_name, is_demo, created_at
		FROM users
		WHERE id = $1
		LIMIT 1`

	var u domain.User
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.IsDemo, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetDemoUser(ctx context.Context) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, display_name, is_demo, created_at
		FROM users
		WHERE is_demo = TRUE
		LIMIT 1`

	var u domain.User
	err := r.db.QueryRow(ctx, q).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.IsDemo, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
