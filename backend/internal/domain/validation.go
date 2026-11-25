package domain

import "time"

func (c IdeaCategory) Valid() bool {
	switch c {
	case CategorySaaS, CategoryMarketplace, CategoryConsumer, CategoryDevtool,
		CategoryHardware, CategoryFintech, CategoryHealth, CategoryEducation,
		CategorySocial, CategoryOther:
		return true
	default:
		return false
	}
}

func (c DeathCause) Valid() bool {
	switch c {
	case CauseNoMarketNeed, CauseRanOutOfMoney, CauseRanOutOfTime,
		CauseLostMotivation, CauseCompetitor, CauseTechnicalDifficulty,
		CauseCofounderConflict, CauseLegalRegulatory, CauseBadTiming,
		CausePivotedAway, CauseOther:
		return true
	default:
		return false
	}
}

func (v ValidationLevel) Valid() bool {
	switch v {
	case ValidationNone, ValidationTalkedToUsers, ValidationPrototype, ValidationPayingCustomers:
		return true
	default:
		return false
	}
}

func (s IdeaStatus) Valid() bool {
	switch s {
	case StatusBuried, StatusOnLifeSupport, StatusResurrected:
		return true
	default:
		return false
	}
}

type IdeaFilter struct {
	Page       int
	Limit      int
	Cause      *DeathCause
	Category   *IdeaCategory
	Status     *IdeaStatus
	Validation *ValidationLevel
	From       *time.Time
	To         *time.Time
	Query      string
	Sort       string
	Order      string
}

type PaginatedIdeas struct {
	Items   []Idea `json:"items"`
	Total   int64  `json:"total"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	HasMore bool   `json:"has_more"`
}
