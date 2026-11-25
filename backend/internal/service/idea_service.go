package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/repository"
)

type CreateIdeaInput struct {
	Name            string
	Tagline         string
	Description     string
	Category        domain.IdeaCategory
	BornAt          time.Time
	DiedAt          time.Time
	PrimaryCause    domain.DeathCause
	SecondaryCauses []domain.DeathCause
	Epitaph         string
	Autopsy         string
	ValidationLevel domain.ValidationLevel
	MoneySpentCents int64
	HoursInvested   int32
	WouldResurrect  bool
}

type UpdateIdeaInput struct {
	Name            *string
	Tagline         *string
	Description     *string
	Category        *domain.IdeaCategory
	BornAt          *time.Time
	DiedAt          *time.Time
	PrimaryCause    *domain.DeathCause
	SecondaryCauses *[]domain.DeathCause
	Epitaph         *string
	Autopsy         *string
	ValidationLevel *domain.ValidationLevel
	MoneySpentCents *int64
	HoursInvested   *int32
	Status          *domain.IdeaStatus
	WouldResurrect  *bool
}

type IdeaService struct {
	ideas        *repository.IdeaRepository
	queryTimeout time.Duration
}

func NewIdeaService(pool *pgxpool.Pool, queryTimeout time.Duration) *IdeaService {
	return &IdeaService{
		ideas:        repository.NewIdeaRepository(pool),
		queryTimeout: queryTimeout,
	}
}

func (s *IdeaService) Create(ctx context.Context, userID uuid.UUID, in CreateIdeaInput) (*domain.Idea, error) {
	if err := ValidateCreateIdea(in); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	secondary := in.SecondaryCauses
	if secondary == nil {
		secondary = []domain.DeathCause{}
	}
	validation := in.ValidationLevel
	if validation == "" {
		validation = domain.ValidationNone
	}

	idea := &domain.Idea{
		UserID:          userID,
		Name:            strings.TrimSpace(in.Name),
		Tagline:         strings.TrimSpace(in.Tagline),
		Description:     strings.TrimSpace(in.Description),
		Category:        in.Category,
		BornAt:          in.BornAt,
		DiedAt:          in.DiedAt,
		PrimaryCause:    in.PrimaryCause,
		SecondaryCauses: secondary,
		Epitaph:         strings.TrimSpace(in.Epitaph),
		Autopsy:         strings.TrimSpace(in.Autopsy),
		ValidationLevel: validation,
		MoneySpentCents: in.MoneySpentCents,
		HoursInvested:   in.HoursInvested,
		Status:          domain.StatusBuried,
		WouldResurrect:  in.WouldResurrect,
	}

	return s.ideas.Create(ctx, idea)
}

func (s *IdeaService) Get(ctx context.Context, userID, ideaID uuid.UUID) (*domain.Idea, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()
	return s.ideas.GetByID(ctx, userID, ideaID)
}

func (s *IdeaService) List(ctx context.Context, userID uuid.UUID, f domain.IdeaFilter) (*domain.PaginatedIdeas, error) {
	f = NormalizeIdeaFilter(f)

	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	items, total, err := s.ideas.List(ctx, userID, f)
	if err != nil {
		return nil, err
	}

	return &domain.PaginatedIdeas{
		Items:   items,
		Total:   total,
		Page:    f.Page,
		Limit:   f.Limit,
		HasMore: int64(f.Page*f.Limit) < total,
	}, nil
}

func (s *IdeaService) Update(ctx context.Context, userID, ideaID uuid.UUID, in UpdateIdeaInput) (*domain.Idea, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	existing, err := s.ideas.GetByID(ctx, userID, ideaID)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		existing.Name = strings.TrimSpace(*in.Name)
	}
	if in.Tagline != nil {
		existing.Tagline = strings.TrimSpace(*in.Tagline)
	}
	if in.Description != nil {
		existing.Description = strings.TrimSpace(*in.Description)
	}
	if in.Category != nil {
		if !in.Category.Valid() {
			return nil, fmt.Errorf("%w: invalid category", domain.ErrValidation)
		}
		existing.Category = *in.Category
	}
	if in.BornAt != nil {
		existing.BornAt = *in.BornAt
	}
	if in.DiedAt != nil {
		existing.DiedAt = *in.DiedAt
	}
	if existing.DiedAt.Before(existing.BornAt) {
		return nil, fmt.Errorf("%w: died_at must be on or after born_at", domain.ErrValidation)
	}
	if in.PrimaryCause != nil {
		if !in.PrimaryCause.Valid() {
			return nil, fmt.Errorf("%w: invalid primary_cause", domain.ErrValidation)
		}
		existing.PrimaryCause = *in.PrimaryCause
	}
	if in.SecondaryCauses != nil {
		for _, c := range *in.SecondaryCauses {
			if !c.Valid() {
				return nil, fmt.Errorf("%w: invalid secondary_cause", domain.ErrValidation)
			}
		}
		existing.SecondaryCauses = *in.SecondaryCauses
	}
	if in.Epitaph != nil {
		trimmed := strings.TrimSpace(*in.Epitaph)
		if trimmed == "" || utf8.RuneCountInString(trimmed) > 140 {
			return nil, fmt.Errorf("%w: epitaph must be between 1 and 140 characters", domain.ErrValidation)
		}
		existing.Epitaph = trimmed
	}
	if in.Autopsy != nil {
		existing.Autopsy = strings.TrimSpace(*in.Autopsy)
	}
	if in.ValidationLevel != nil {
		if !in.ValidationLevel.Valid() {
			return nil, fmt.Errorf("%w: invalid validation_level", domain.ErrValidation)
		}
		existing.ValidationLevel = *in.ValidationLevel
	}
	if in.MoneySpentCents != nil {
		if *in.MoneySpentCents < 0 {
			return nil, fmt.Errorf("%w: money_spent_cents cannot be negative", domain.ErrValidation)
		}
		existing.MoneySpentCents = *in.MoneySpentCents
	}
	if in.HoursInvested != nil {
		if *in.HoursInvested < 0 {
			return nil, fmt.Errorf("%w: hours_invested cannot be negative", domain.ErrValidation)
		}
		existing.HoursInvested = *in.HoursInvested
	}
	if in.Status != nil && *in.Status != existing.Status {
		if !existing.Status.CanTransitionTo(*in.Status) {
			return nil, domain.ErrInvalidTransition
		}
		existing.Status = *in.Status
	}
	if in.WouldResurrect != nil {
		existing.WouldResurrect = *in.WouldResurrect
	}

	return s.ideas.Update(ctx, existing)
}

func (s *IdeaService) Delete(ctx context.Context, userID, ideaID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()
	return s.ideas.SoftDelete(ctx, userID, ideaID)
}

func (s *IdeaService) Resurrect(ctx context.Context, userID, ideaID uuid.UUID) (*domain.Idea, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	existing, err := s.ideas.GetByID(ctx, userID, ideaID)
	if err != nil {
		return nil, err
	}

	next, ok := existing.Status.NextResurrectionStatus()
	if !ok {
		return nil, domain.ErrInvalidTransition
	}

	return s.ideas.UpdateStatus(ctx, userID, ideaID, next)
}

func ValidateCreateIdea(in CreateIdeaInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrValidation)
	}
	if !in.Category.Valid() {
		return fmt.Errorf("%w: invalid category", domain.ErrValidation)
	}
	if in.DiedAt.Before(in.BornAt) {
		return fmt.Errorf("%w: died_at must be on or after born_at", domain.ErrValidation)
	}
	if !in.PrimaryCause.Valid() {
		return fmt.Errorf("%w: invalid primary_cause", domain.ErrValidation)
	}
	for _, c := range in.SecondaryCauses {
		if !c.Valid() {
			return fmt.Errorf("%w: invalid secondary_cause", domain.ErrValidation)
		}
	}
	epitaph := strings.TrimSpace(in.Epitaph)
	if epitaph == "" || utf8.RuneCountInString(epitaph) > 140 {
		return fmt.Errorf("%w: epitaph must be between 1 and 140 characters", domain.ErrValidation)
	}
	if in.ValidationLevel != "" && !in.ValidationLevel.Valid() {
		return fmt.Errorf("%w: invalid validation_level", domain.ErrValidation)
	}
	if in.MoneySpentCents < 0 {
		return fmt.Errorf("%w: money_spent_cents cannot be negative", domain.ErrValidation)
	}
	if in.HoursInvested < 0 {
		return fmt.Errorf("%w: hours_invested cannot be negative", domain.ErrValidation)
	}
	return nil
}

func NormalizeIdeaFilter(f domain.IdeaFilter) domain.IdeaFilter {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 12
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Sort == "" {
		f.Sort = "died_at"
	}
	if f.Order != "asc" && f.Order != "desc" {
		f.Order = "desc"
	}
	return f
}
