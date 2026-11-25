package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nisha/gravehard/internal/auth"
	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/repository"
)

type AuthResult struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	pool         *pgxpool.Pool
	users        *repository.UserRepository
	tokens       *repository.RefreshTokenRepository
	jwtMgr       *auth.JWTManager
	refreshTTL   time.Duration
	queryTimeout time.Duration
}

func NewAuthService(
	pool *pgxpool.Pool,
	jwtMgr *auth.JWTManager,
	refreshTTL time.Duration,
	queryTimeout time.Duration,
) *AuthService {
	return &AuthService{
		pool:         pool,
		users:        repository.NewUserRepository(pool),
		tokens:       repository.NewRefreshTokenRepository(pool),
		jwtMgr:       jwtMgr,
		refreshTTL:   refreshTTL,
		queryTimeout: queryTimeout,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (*AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	email = strings.TrimSpace(strings.ToLower(email))
	displayName = strings.TrimSpace(displayName)

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, email, hash, displayName, false)
	if err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user, uuid.New())
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, err
	}

	ok, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !ok {
		return nil, domain.ErrUnauthorized
	}

	return s.issueSession(ctx, user, uuid.New())
}

// Refresh validates the refresh token, rotates it inside a transaction, and revokes the family on reuse.
func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (*AuthResult, error) {
	if rawRefreshToken == "" {
		return nil, domain.ErrUnauthorized
	}

	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	tokenHash := auth.HashRefreshToken(rawRefreshToken)
	rec, err := s.tokens.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, err
	}

	// Reuse detection: if an already-revoked token is presented, revoke the entire family.
	if rec.RevokedAt != nil {
		_ = s.tokens.RevokeFamily(ctx, rec.FamilyID)
		return nil, domain.ErrUnauthorized
	}

	if time.Now().UTC().After(rec.ExpiresAt) {
		return nil, domain.ErrUnauthorized
	}

	user, err := s.users.GetByID(ctx, rec.UserID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	newRaw, newHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = repository.WithTx(ctx, s.pool, func(tx pgx.Tx) error {
		txTokens := repository.NewRefreshTokenRepository(tx)
		if err := txTokens.RevokeByID(ctx, rec.ID); err != nil {
			return err
		}
		_, err := txTokens.Create(ctx, user.ID, newHash, rec.FamilyID, time.Now().UTC().Add(s.refreshTTL))
		return err
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtMgr.IssueAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRaw,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	tokenHash := auth.HashRefreshToken(rawRefreshToken)
	rec, err := s.tokens.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}
	return s.tokens.RevokeFamily(ctx, rec.FamilyID)
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()
	return s.users.GetByID(ctx, userID)
}

func (s *AuthService) issueSession(ctx context.Context, user *domain.User, familyID uuid.UUID) (*AuthResult, error) {
	accessToken, err := s.jwtMgr.IssueAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	rawRefresh, hashRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	_, err = s.tokens.Create(ctx, user.ID, hashRefresh, familyID, time.Now().UTC().Add(s.refreshTTL))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
