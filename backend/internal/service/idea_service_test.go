package service_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/service"
)

func TestIdeaStatusTransitions(t *testing.T) {
	// buried -> on_life_support -> resurrected
	next, ok := domain.StatusBuried.NextResurrectionStatus()
	require.True(t, ok)
	require.Equal(t, domain.StatusOnLifeSupport, next)

	next, ok = domain.StatusOnLifeSupport.NextResurrectionStatus()
	require.True(t, ok)
	require.Equal(t, domain.StatusResurrected, next)

	_, ok = domain.StatusResurrected.NextResurrectionStatus()
	require.False(t, ok)

	// CanTransitionTo rejects skipping or backwards
	require.True(t, domain.StatusBuried.CanTransitionTo(domain.StatusOnLifeSupport))
	require.False(t, domain.StatusBuried.CanTransitionTo(domain.StatusResurrected))
	require.False(t, domain.StatusResurrected.CanTransitionTo(domain.StatusBuried))
	require.False(t, domain.StatusOnLifeSupport.CanTransitionTo(domain.StatusBuried))
}

func TestValidateCreateIdea(t *testing.T) {
	born := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	died := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	valid := service.CreateIdeaInput{
		Name:            "DogWalkerCRM",
		Category:        domain.CategorySaaS,
		BornAt:          born,
		DiedAt:          died,
		PrimaryCause:    domain.CauseNoMarketNeed,
		Epitaph:         "Here lies a CRM for dog walkers. It never met a customer.",
		ValidationLevel: domain.ValidationNone,
		MoneySpentCents: 12000,
		HoursInvested:   40,
	}
	require.NoError(t, service.ValidateCreateIdea(valid))

	// died_at before born_at
	badDates := valid
	badDates.DiedAt = born.AddDate(0, 0, -1)
	require.ErrorIs(t, service.ValidateCreateIdea(badDates), domain.ErrValidation)

	// epitaph > 140 chars
	longEpitaph := valid
	longEpitaph.Epitaph = strings.Repeat("a", 141)
	require.ErrorIs(t, service.ValidateCreateIdea(longEpitaph), domain.ErrValidation)

	// negative money
	negMoney := valid
	negMoney.MoneySpentCents = -100
	require.ErrorIs(t, service.ValidateCreateIdea(negMoney), domain.ErrValidation)
}

func TestNormalizeIdeaFilter(t *testing.T) {
	f := service.NormalizeIdeaFilter(domain.IdeaFilter{Page: -2, Limit: 500})
	require.Equal(t, 1, f.Page)
	require.Equal(t, 100, f.Limit)
	require.Equal(t, "died_at", f.Sort)
	require.Equal(t, "desc", f.Order)
}
