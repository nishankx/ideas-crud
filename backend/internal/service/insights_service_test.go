package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/service"
)

func makeTestIdea(cat domain.IdeaCategory, cause domain.DeathCause, val domain.ValidationLevel, days int, cents int64, hours int32, resurrect bool) domain.Idea {
	born := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return domain.Idea{
		Category:        cat,
		BornAt:          born,
		DiedAt:          born.AddDate(0, 0, days),
		PrimaryCause:    cause,
		SecondaryCauses: []domain.DeathCause{domain.CauseLostMotivation},
		ValidationLevel: val,
		MoneySpentCents: cents,
		HoursInvested:   hours,
		Status:          domain.StatusBuried,
		WouldResurrect:  resurrect,
	}
}

func TestComputePatternsMinimumSampleSize(t *testing.T) {
	// Fewer than 3 ideas -> must emit 0 patterns
	ideas := []domain.Idea{
		makeTestIdea(domain.CategorySaaS, domain.CauseNoMarketNeed, domain.ValidationNone, 10, 1000, 20, false),
		makeTestIdea(domain.CategorySaaS, domain.CauseNoMarketNeed, domain.ValidationNone, 10, 1000, 20, false),
	}
	require.Empty(t, service.ComputePatterns(ideas))
}

func TestComputePatternsDetectsInsights(t *testing.T) {
	ideas := []domain.Idea{
		makeTestIdea(domain.CategoryFintech, domain.CauseNoMarketNeed, domain.ValidationNone, 5, 10000, 30, true),
		makeTestIdea(domain.CategoryFintech, domain.CauseNoMarketNeed, domain.ValidationNone, 6, 20000, 40, true),
		makeTestIdea(domain.CategoryFintech, domain.CauseNoMarketNeed, domain.ValidationNone, 7, 15000, 50, true),
		makeTestIdea(domain.CategorySaaS, domain.CauseLostMotivation, domain.ValidationPrototype, 90, 5000, 200, false),
		makeTestIdea(domain.CategoryDevtool, domain.CauseRanOutOfTime, domain.ValidationPayingCustomers, 120, 8000, 180, false),
		makeTestIdea(domain.CategoryConsumer, domain.CauseCompetitor, domain.ValidationTalkedToUsers, 100, 3000, 160, false),
	}

	patterns := service.ComputePatterns(ideas)
	ids := make(map[string]bool)
	for _, p := range patterns {
		require.GreaterOrEqual(t, p.SampleSize, 3)
		ids[p.ID] = true
	}

	require.True(t, ids["unvalidated_no_market"])
	require.True(t, ids["fast_dying_category"])
	require.True(t, ids["serial_killer_cause"])
	require.True(t, ids["marathon_burnout"])
	require.True(t, ids["unfinished_business"])
}

func TestComputeSummaryAndBuckets(t *testing.T) {
	ideas := []domain.Idea{
		makeTestIdea(domain.CategorySaaS, domain.CauseNoMarketNeed, domain.ValidationNone, 3, 1000, 10, false),
		makeTestIdea(domain.CategorySaaS, domain.CauseNoMarketNeed, domain.ValidationTalkedToUsers, 14, 2000, 20, false),
		makeTestIdea(domain.CategoryFintech, domain.CauseCompetitor, domain.ValidationPrototype, 45, 3000, 30, false),
		makeTestIdea(domain.CategoryHealth, domain.CauseLegalRegulatory, domain.ValidationPayingCustomers, 120, 4000, 40, false),
		makeTestIdea(domain.CategoryOther, domain.CauseBadTiming, domain.ValidationNone, 200, 5000, 50, false),
	}

	sum := service.ComputeSummary(ideas)
	require.Equal(t, 5, sum.TotalIdeas)
	require.Equal(t, int64(15000), sum.TotalMoneyLostCents)
	require.Equal(t, int64(150), sum.TotalHoursLost)
	require.NotNil(t, sum.MostCommonCause)
	require.Equal(t, domain.CauseNoMarketNeed, *sum.MostCommonCause)
	require.NotNil(t, sum.MostLethalCategory)
	require.Equal(t, domain.CategorySaaS, *sum.MostLethalCategory)

	lifespan := service.ComputeLifespan(ideas)
	require.Len(t, lifespan.Buckets, 5)
	for _, b := range lifespan.Buckets {
		require.Equal(t, 1, b.Count)
	}

	causes := service.ComputeCauses(ideas)
	require.Equal(t, 5, causes.Total)
	require.Equal(t, domain.CauseNoMarketNeed, causes.Causes[0].Cause)
	require.Equal(t, 2, causes.Causes[0].Count)

	val := service.ComputeValidation(ideas)
	require.Len(t, val.Matrix, 4)
}
