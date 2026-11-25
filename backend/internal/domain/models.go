package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdeaCategory string

const (
	CategorySaaS        IdeaCategory = "saas"
	CategoryMarketplace IdeaCategory = "marketplace"
	CategoryConsumer    IdeaCategory = "consumer"
	CategoryDevtool     IdeaCategory = "devtool"
	CategoryHardware    IdeaCategory = "hardware"
	CategoryFintech     IdeaCategory = "fintech"
	CategoryHealth      IdeaCategory = "health"
	CategoryEducation   IdeaCategory = "education"
	CategorySocial      IdeaCategory = "social"
	CategoryOther       IdeaCategory = "other"
)

type DeathCause string

const (
	CauseNoMarketNeed        DeathCause = "no_market_need"
	CauseRanOutOfMoney       DeathCause = "ran_out_of_money"
	CauseRanOutOfTime        DeathCause = "ran_out_of_time"
	CauseLostMotivation      DeathCause = "lost_motivation"
	CauseCompetitor          DeathCause = "competitor"
	CauseTechnicalDifficulty DeathCause = "technical_difficulty"
	CauseCofounderConflict   DeathCause = "cofounder_conflict"
	CauseLegalRegulatory     DeathCause = "legal_regulatory"
	CauseBadTiming           DeathCause = "bad_timing"
	CausePivotedAway         DeathCause = "pivoted_away"
	CauseOther               DeathCause = "other"
)

type ValidationLevel string

const (
	ValidationNone            ValidationLevel = "none"
	ValidationTalkedToUsers   ValidationLevel = "talked_to_users"
	ValidationPrototype       ValidationLevel = "prototype"
	ValidationPayingCustomers ValidationLevel = "paying_customers"
)

type IdeaStatus string

const (
	StatusBuried        IdeaStatus = "buried"
	StatusOnLifeSupport IdeaStatus = "on_life_support"
	StatusResurrected   IdeaStatus = "resurrected"
)

// CanTransitionTo enforces valid idea resurrection state transitions:
// buried -> on_life_support -> resurrected.
func (s IdeaStatus) CanTransitionTo(next IdeaStatus) bool {
	switch s {
	case StatusBuried:
		return next == StatusOnLifeSupport
	case StatusOnLifeSupport:
		return next == StatusResurrected
	default:
		return false
	}
}

// NextResurrectionStatus returns the next status in the resurrection pipeline.
func (s IdeaStatus) NextResurrectionStatus() (IdeaStatus, bool) {
	switch s {
	case StatusBuried:
		return StatusOnLifeSupport, true
	case StatusOnLifeSupport:
		return StatusResurrected, true
	default:
		return "", false
	}
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	IsDemo       bool      `json:"is_demo"`
	CreatedAt    time.Time `json:"created_at"`
}

type Idea struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	Name            string          `json:"name"`
	Tagline         string          `json:"tagline"`
	Description     string          `json:"description"`
	Category        IdeaCategory    `json:"category"`
	BornAt          time.Time       `json:"born_at"`
	DiedAt          time.Time       `json:"died_at"`
	PrimaryCause    DeathCause      `json:"primary_cause"`
	SecondaryCauses []DeathCause    `json:"secondary_causes"`
	Epitaph         string          `json:"epitaph"`
	Autopsy         string          `json:"autopsy"`
	ValidationLevel ValidationLevel `json:"validation_level"`
	MoneySpentCents int64           `json:"money_spent_cents"`
	HoursInvested   int32           `json:"hours_invested"`
	Status          IdeaStatus      `json:"status"`
	WouldResurrect  bool            `json:"would_resurrect"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
