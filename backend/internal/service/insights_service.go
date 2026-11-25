package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/repository"
)

const minPatternSampleSize = 3

type InsightsService struct {
	ideas        *repository.IdeaRepository
	queryTimeout time.Duration
}

func NewInsightsService(pool *pgxpool.Pool, queryTimeout time.Duration) *InsightsService {
	return &InsightsService{
		ideas:        repository.NewIdeaRepository(pool),
		queryTimeout: queryTimeout,
	}
}

func (s *InsightsService) Summary(ctx context.Context, userID uuid.UUID) (*domain.InsightsSummary, error) {
	ideas, err := s.fetchAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ComputeSummary(ideas), nil
}

func (s *InsightsService) Causes(ctx context.Context, userID uuid.UUID) (*domain.CausesResponse, error) {
	ideas, err := s.fetchAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ComputeCauses(ideas), nil
}

func (s *InsightsService) Timeline(ctx context.Context, userID uuid.UUID) (*domain.TimelineResponse, error) {
	ideas, err := s.fetchAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ComputeTimeline(ideas), nil
}

func (s *InsightsService) Lifespan(ctx context.Context, userID uuid.UUID) (*domain.LifespanResponse, error) {
	ideas, err := s.fetchAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ComputeLifespan(ideas), nil
}

func (s *InsightsService) Validation(ctx context.Context, userID uuid.UUID) (*domain.ValidationResponse, error) {
	ideas, err := s.fetchAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ComputeValidation(ideas), nil
}

func (s *InsightsService) Patterns(ctx context.Context, userID uuid.UUID) (*domain.PatternsResponse, error) {
	ideas, err := s.fetchAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &domain.PatternsResponse{Patterns: ComputePatterns(ideas)}, nil
}

func (s *InsightsService) fetchAll(ctx context.Context, userID uuid.UUID) ([]domain.Idea, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()
	return s.ideas.ListAllByUser(ctx, userID)
}

func lifespanDays(idea domain.Idea) float64 {
	days := idea.DiedAt.Sub(idea.BornAt).Hours() / 24.0
	if days < 0 {
		return 0
	}
	return days
}

func ComputeSummary(ideas []domain.Idea) *domain.InsightsSummary {
	n := len(ideas)
	if n == 0 {
		return &domain.InsightsSummary{}
	}

	var money int64
	var hours int64
	var totalDays float64
	causeCounts := make(map[domain.DeathCause]int)
	catCounts := make(map[domain.IdeaCategory]int)

	for _, it := range ideas {
		money += it.MoneySpentCents
		hours += int64(it.HoursInvested)
		totalDays += lifespanDays(it)
		causeCounts[it.PrimaryCause]++
		catCounts[it.Category]++
	}

	var topCause domain.DeathCause
	topCauseCount := -1
	for c, cnt := range causeCounts {
		if cnt > topCauseCount || (cnt == topCauseCount && c < topCause) {
			topCause = c
			topCauseCount = cnt
		}
	}

	var topCat domain.IdeaCategory
	topCatCount := -1
	for cat, cnt := range catCounts {
		if cnt > topCatCount || (cnt == topCatCount && cat < topCat) {
			topCat = cat
			topCatCount = cnt
		}
	}

	return &domain.InsightsSummary{
		TotalIdeas:          n,
		TotalMoneyLostCents: money,
		TotalHoursLost:      hours,
		AvgLifespanDays:     math.Round((totalDays/float64(n))*10) / 10,
		MostCommonCause:     &topCause,
		MostLethalCategory:  &topCat,
	}
}

func ComputeCauses(ideas []domain.Idea) *domain.CausesResponse {
	n := len(ideas)
	counts := make(map[domain.DeathCause]int)
	coMap := make(map[domain.DeathCause]map[domain.DeathCause]int)

	for _, it := range ideas {
		counts[it.PrimaryCause]++
		if coMap[it.PrimaryCause] == nil {
			coMap[it.PrimaryCause] = make(map[domain.DeathCause]int)
		}
		for _, sec := range it.SecondaryCauses {
			if sec != it.PrimaryCause {
				coMap[it.PrimaryCause][sec]++
			}
		}
	}

	items := make([]domain.CauseBreakdownItem, 0, len(counts))
	for cause, cnt := range counts {
		coList := make([]domain.CoOccurringCause, 0, len(coMap[cause]))
		for sc, scCnt := range coMap[cause] {
			coList = append(coList, domain.CoOccurringCause{Cause: sc, Count: scCnt})
		}
		sort.Slice(coList, func(i, j int) bool {
			if coList[i].Count != coList[j].Count {
				return coList[i].Count > coList[j].Count
			}
			return coList[i].Cause < coList[j].Cause
		})

		share := 0.0
		if n > 0 {
			share = math.Round((float64(cnt)/float64(n))*1000) / 1000
		}
		items = append(items, domain.CauseBreakdownItem{
			Cause:       cause,
			Count:       cnt,
			Share:       share,
			CoOccurring: coList,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Cause < items[j].Cause
	})

	return &domain.CausesResponse{Total: n, Causes: items}
}

func ComputeTimeline(ideas []domain.Idea) *domain.TimelineResponse {
	monthCounts := make(map[string]int)
	for _, it := range ideas {
		key := it.DiedAt.Format("2006-01")
		monthCounts[key]++
	}

	months := make([]string, 0, len(monthCounts))
	for m := range monthCounts {
		months = append(months, m)
	}
	sort.Strings(months)

	buckets := make([]domain.TimelineBucket, 0, len(months))
	for _, m := range months {
		buckets = append(buckets, domain.TimelineBucket{Month: m, Count: monthCounts[m]})
	}
	return &domain.TimelineResponse{Buckets: buckets}
}

func ComputeLifespan(ideas []domain.Idea) *domain.LifespanResponse {
	defs := []struct {
		key   string
		label string
	}{
		{"<1w", "< 1 week"},
		{"1-4w", "1–4 weeks"},
		{"1-3m", "1–3 months"},
		{"3-6m", "3–6 months"},
		{"6m+", "6+ months"},
	}

	counts := make([]int, len(defs))
	for _, it := range ideas {
		days := lifespanDays(it)
		switch {
		case days < 7:
			counts[0]++
		case days < 28:
			counts[1]++
		case days < 90:
			counts[2]++
		case days < 180:
			counts[3]++
		default:
			counts[4]++
		}
	}

	n := len(ideas)
	buckets := make([]domain.LifespanBucket, len(defs))
	for i, d := range defs {
		share := 0.0
		if n > 0 {
			share = math.Round((float64(counts[i])/float64(n))*1000) / 1000
		}
		buckets[i] = domain.LifespanBucket{
			Bucket: d.key,
			Label:  d.label,
			Count:  counts[i],
			Share:  share,
		}
	}
	return &domain.LifespanResponse{Total: n, Buckets: buckets}
}

func ComputeValidation(ideas []domain.Idea) *domain.ValidationResponse {
	levels := []domain.ValidationLevel{
		domain.ValidationNone,
		domain.ValidationTalkedToUsers,
		domain.ValidationPrototype,
		domain.ValidationPayingCustomers,
	}

	byLevel := make(map[domain.ValidationLevel]*domain.ValidationMatrixRow, len(levels))
	for _, lvl := range levels {
		byLevel[lvl] = &domain.ValidationMatrixRow{
			ValidationLevel: lvl,
			Causes:          make(map[domain.DeathCause]int),
		}
	}

	for _, it := range ideas {
		row, ok := byLevel[it.ValidationLevel]
		if !ok {
			continue
		}
		row.Total++
		row.Causes[it.PrimaryCause]++
	}

	matrix := make([]domain.ValidationMatrixRow, 0, len(levels))
	for _, lvl := range levels {
		matrix = append(matrix, *byLevel[lvl])
	}
	return &domain.ValidationResponse{Matrix: matrix}
}

// ComputePatterns generates plain-English observations when sample size >= 3.
func ComputePatterns(ideas []domain.Idea) []domain.PatternObservation {
	patterns := make([]domain.PatternObservation, 0)
	n := len(ideas)
	if n < minPatternSampleSize {
		return patterns
	}

	// 1. Unvalidated ideas dying of no_market_need
	var unvalidatedTotal, unvalidatedNoMarket int
	for _, it := range ideas {
		if it.ValidationLevel == domain.ValidationNone {
			unvalidatedTotal++
			if it.PrimaryCause == domain.CauseNoMarketNeed {
				unvalidatedNoMarket++
			}
		}
	}
	if unvalidatedTotal >= minPatternSampleSize {
		pct := int(math.Round(float64(unvalidatedNoMarket) * 100 / float64(unvalidatedTotal)))
		if pct >= 45 {
			patterns = append(patterns, domain.PatternObservation{
				ID:         "unvalidated_no_market",
				Title:      "Building in a vacuum",
				Message:    fmt.Sprintf("%d%% of your ideas that never reached users died of No Market Need.", pct),
				Severity:   "critical",
				SampleSize: unvalidatedTotal,
			})
		}
	}

	// 2. Category dying significantly faster than overall average
	var totalDays float64
	catDays := make(map[domain.IdeaCategory]float64)
	catCount := make(map[domain.IdeaCategory]int)
	for _, it := range ideas {
		d := lifespanDays(it)
		totalDays += d
		catDays[it.Category] += d
		catCount[it.Category]++
	}
	overallAvg := totalDays / float64(n)
	if overallAvg > 0 {
		type fastCat struct {
			cat   domain.IdeaCategory
			ratio float64
			cnt   int
		}
		var fastest *fastCat
		for cat, cnt := range catCount {
			if cnt < minPatternSampleSize {
				continue
			}
			catAvg := catDays[cat] / float64(cnt)
			if catAvg > 0 && overallAvg/catAvg >= 1.75 {
				ratio := math.Round((overallAvg/catAvg)*10) / 10
				if fastest == nil || ratio > fastest.ratio {
					fastest = &fastCat{cat: cat, ratio: ratio, cnt: cnt}
				}
			}
		}
		if fastest != nil {
			patterns = append(patterns, domain.PatternObservation{
				ID:         "fast_dying_category",
				Title:      "Short shelf life",
				Message:    fmt.Sprintf("Your %s ideas die %.1fx faster than your overall average.", fastest.cat, fastest.ratio),
				Severity:   "warning",
				SampleSize: fastest.cnt,
			})
		}
	}

	// 3. Dominant primary cause (>= 35% of all ideas and >= 3)
	causeCounts := make(map[domain.DeathCause]int)
	for _, it := range ideas {
		causeCounts[it.PrimaryCause]++
	}
	var topCause domain.DeathCause
	var topCount int
	for c, cnt := range causeCounts {
		if cnt > topCount {
			topCause = c
			topCount = cnt
		}
	}
	if topCount >= minPatternSampleSize {
		pct := int(math.Round(float64(topCount) * 100 / float64(n)))
		if pct >= 35 {
			patterns = append(patterns, domain.PatternObservation{
				ID:         "serial_killer_cause",
				Title:      "Repeat offender",
				Message:    fmt.Sprintf("%s is your #1 killer — responsible for %d%% of all burials.", humanizeCause(topCause), pct),
				Severity:   "warning",
				SampleSize: n,
			})
		}
	}

	// 4. High-hour burnout (150+ hours ending in lost_motivation or ran_out_of_time)
	var heavyTotal, heavyBurnout int
	for _, it := range ideas {
		if it.HoursInvested >= 150 {
			heavyTotal++
			if it.PrimaryCause == domain.CauseLostMotivation || it.PrimaryCause == domain.CauseRanOutOfTime {
				heavyBurnout++
			}
		}
	}
	if heavyTotal >= minPatternSampleSize && heavyBurnout >= 2 {
		patterns = append(patterns, domain.PatternObservation{
			ID:         "marathon_burnout",
			Title:      "Marathon fatigue",
			Message:    fmt.Sprintf("%d of your %d ideas with 150+ hours invested ended in Lost Motivation or Running Out of Time.", heavyBurnout, heavyTotal),
			Severity:   "info",
			SampleSize: heavyTotal,
		})
	}

	// 5. Resurrection candidates
	var resurrectable int
	for _, it := range ideas {
		if it.WouldResurrect && it.Status == domain.StatusBuried {
			resurrectable++
		}
	}
	if resurrectable >= minPatternSampleSize {
		patterns = append(patterns, domain.PatternObservation{
			ID:         "unfinished_business",
			Title:      "Unfinished business",
			Message:    fmt.Sprintf("You marked %d buried ideas as worth resurrecting — consider reviving one before starting from scratch.", resurrectable),
			Severity:   "info",
			SampleSize: resurrectable,
		})
	}

	return patterns
}

func humanizeCause(c domain.DeathCause) string {
	parts := strings.Split(string(c), "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
