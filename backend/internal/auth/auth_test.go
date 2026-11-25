package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/nisha/gravehard/internal/auth"
)

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := auth.HashPassword("super-secret-password")
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	ok, err := auth.VerifyPassword("super-secret-password", hash)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = auth.VerifyPassword("wrong-password", hash)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestJWTIssueAndVerify(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret-key", 15*time.Minute)
	uid := uuid.New()

	tok, err := mgr.IssueAccessToken(uid, "founder@gravehard.dev")
	require.NoError(t, err)

	claims, err := mgr.VerifyAccessToken(tok)
	require.NoError(t, err)
	require.Equal(t, uid, claims.UserID)
	require.Equal(t, "founder@gravehard.dev", claims.Email)

	otherMgr := auth.NewJWTManager("wrong-secret", 15*time.Minute)
	_, err = otherMgr.VerifyAccessToken(tok)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestJWTExpired(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret-key", -1*time.Minute)
	tok, err := mgr.IssueAccessToken(uuid.New(), "expired@gravehard.dev")
	require.NoError(t, err)

	_, err = mgr.VerifyAccessToken(tok)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestRefreshTokenGeneration(t *testing.T) {
	raw1, hash1, err := auth.GenerateRefreshToken()
	require.NoError(t, err)
	raw2, hash2, err := auth.GenerateRefreshToken()
	require.NoError(t, err)

	require.NotEqual(t, raw1, raw2)
	require.NotEqual(t, hash1, hash2)
	require.Equal(t, hash1, auth.HashRefreshToken(raw1))
}
