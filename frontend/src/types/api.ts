export type IdeaCategory =
  | 'saas'
  | 'marketplace'
  | 'consumer'
  | 'devtool'
  | 'hardware'
  | 'fintech'
  | 'health'
  | 'education'
  | 'social'
  | 'other';

export type DeathCause =
  | 'no_market_need'
  | 'ran_out_of_money'
  | 'ran_out_of_time'
  | 'lost_motivation'
  | 'competitor'
  | 'technical_difficulty'
  | 'cofounder_conflict'
  | 'legal_regulatory'
  | 'bad_timing'
  | 'pivoted_away'
  | 'other';

export type ValidationLevel =
  | 'none'
  | 'talked_to_users'
  | 'prototype'
  | 'paying_customers';

export type IdeaStatus = 'buried' | 'on_life_support' | 'resurrected';

export interface User {
  id: string;
  email: string;
  display_name: string;
  is_demo: boolean;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  user: User;
}

export interface Idea {
  id: string;
  user_id: string;
  name: string;
  tagline: string;
  description: string;
  category: IdeaCategory;
  born_at: string;
  died_at: string;
  primary_cause: DeathCause;
  secondary_causes: DeathCause[];
  epitaph: string;
  autopsy: string;
  validation_level: ValidationLevel;
  money_spent_cents: number;
  hours_invested: number;
  status: IdeaStatus;
  would_resurrect: boolean;
  created_at: string;
  updated_at: string;
}

export interface PaginatedIdeas {
  items: Idea[];
  total: number;
  page: number;
  limit: number;
  has_more: boolean;
}

export interface IdeaFilterParams {
  page?: number;
  limit?: number;
  cause?: DeathCause | '';
  category?: IdeaCategory | '';
  status?: IdeaStatus | '';
  validation?: ValidationLevel | '';
  q?: string;
  sort?: 'died_at' | 'created_at' | 'money_spent' | 'lifespan';
  order?: 'asc' | 'desc';
}

export interface InsightsSummary {
  total_ideas: number;
  total_money_lost_cents: number;
  total_hours_lost: number;
  avg_lifespan_days: number;
  most_common_cause: DeathCause | null;
  most_lethal_category: IdeaCategory | null;
}

export interface CauseBreakdownItem {
  cause: DeathCause;
  count: number;
  share: number;
  co_occurring: { cause: DeathCause; count: number }[];
}

export interface CausesResponse {
  total: number;
  causes: CauseBreakdownItem[];
}

export interface TimelineBucket {
  month: string;
  count: number;
}

export interface TimelineResponse {
  buckets: TimelineBucket[];
}

export interface LifespanBucket {
  bucket: string;
  label: string;
  count: number;
  share: number;
}

export interface LifespanResponse {
  total: number;
  buckets: LifespanBucket[];
}

export interface ValidationMatrixRow {
  validation_level: ValidationLevel;
  total: number;
  causes: Record<string, number>;
}

export interface ValidationResponse {
  matrix: ValidationMatrixRow[];
}

export interface PatternObservation {
  id: string;
  title: string;
  message: string;
  severity: 'info' | 'warning' | 'critical';
  sample_size: number;
}

export interface PatternsResponse {
  patterns: PatternObservation[];
}

export interface ApiErrorBody {
  code: string;
  message: string;
  details?: Record<string, string>;
}
