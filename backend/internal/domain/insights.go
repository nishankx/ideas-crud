package domain

type InsightsSummary struct {
	TotalIdeas          int           `json:"total_ideas"`
	TotalMoneyLostCents int64         `json:"total_money_lost_cents"`
	TotalHoursLost      int64         `json:"total_hours_lost"`
	AvgLifespanDays     float64       `json:"avg_lifespan_days"`
	MostCommonCause     *DeathCause   `json:"most_common_cause"`
	MostLethalCategory  *IdeaCategory `json:"most_lethal_category"`
}

type CoOccurringCause struct {
	Cause DeathCause `json:"cause"`
	Count int        `json:"count"`
}

type CauseBreakdownItem struct {
	Cause       DeathCause         `json:"cause"`
	Count       int                `json:"count"`
	Share       float64            `json:"share"`
	CoOccurring []CoOccurringCause `json:"co_occurring"`
}

type CausesResponse struct {
	Total  int                  `json:"total"`
	Causes []CauseBreakdownItem `json:"causes"`
}

type TimelineBucket struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type TimelineResponse struct {
	Buckets []TimelineBucket `json:"buckets"`
}

type LifespanBucket struct {
	Bucket string  `json:"bucket"`
	Label  string  `json:"label"`
	Count  int     `json:"count"`
	Share  float64 `json:"share"`
}

type LifespanResponse struct {
	Total   int              `json:"total"`
	Buckets []LifespanBucket `json:"buckets"`
}

type ValidationMatrixRow struct {
	ValidationLevel ValidationLevel    `json:"validation_level"`
	Total           int                `json:"total"`
	Causes          map[DeathCause]int `json:"causes"`
}

type ValidationResponse struct {
	Matrix []ValidationMatrixRow `json:"matrix"`
}

type PatternObservation struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	Severity   string `json:"severity"` // info | warning | critical
	SampleSize int    `json:"sample_size"`
}

type PatternsResponse struct {
	Patterns []PatternObservation `json:"patterns"`
}
