package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/nisha/gravehard/internal/domain"
)

type RefreshTokenRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	FamilyID  uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type RefreshTokenRepository struct {
	db DBTX
}

func NewRefreshTokenRepository(db DBTX) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, familyID uuid.UUID, expiresAt time.Time) (*RefreshTokenRecord, error) {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, family_id, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, token_hash, family_id, expires_at, revoked_at, created_at`

	var rec RefreshTokenRecord
	err := r.db.QueryRow(ctx, q, userID, tokenHash, familyID, expiresAt).Scan(
		&rec.ID, &rec.UserID, &rec.TokenHash, &rec.FamilyID, &rec.ExpiresAt, &rec.RevokedAt, &rec.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RefreshTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error) {
	const q = `
		SELECT id, user_id, token_hash, family_id, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		LIMIT 1`

	var rec RefreshTokenRecord
	err := r.db.QueryRow(ctx, q, tokenHash).Scan(
		&rec.ID, &rec.UserID, &rec.TokenHash, &rec.FamilyID, &rec.ExpiresAt, &rec.RevokedAt, &rec.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RefreshTokenRepository) RevokeByID(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

func (r *RefreshTokenRepository) RevokeFamily(ctx context.Context, familyID uuid.UUID) error {
	const q = `UPDATE refresh_tokens SET revoked_at = NOW() WHERE family_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, q, familyID)
	return err
}
