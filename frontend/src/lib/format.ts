import type { DeathCause, IdeaCategory, IdeaStatus, ValidationLevel } from '../types/api';

export const CAUSE_LABELS: Record<DeathCause, string> = {
  no_market_need: 'No Market Need',
  ran_out_of_money: 'Ran Out of Money',
  ran_out_of_time: 'Ran Out of Time',
  lost_motivation: 'Lost Motivation',
  competitor: 'Crushed by Competitor',
  technical_difficulty: 'Technical Wall',
  cofounder_conflict: 'Cofounder Conflict',
  legal_regulatory: 'Legal / Regulatory',
  bad_timing: 'Bad Timing',
  pivoted_away: 'Pivoted Away',
  other: 'Other',
};

export const CATEGORY_LABELS: Record<IdeaCategory, string> = {
  saas: 'SaaS',
  marketplace: 'Marketplace',
  consumer: 'Consumer',
  devtool: 'DevTool',
  hardware: 'Hardware',
  fintech: 'Fintech',
  health: 'Health',
  education: 'Education',
  social: 'Social',
  other: 'Other',
};

export const VALIDATION_LABELS: Record<ValidationLevel, string> = {
  none: 'Zero Validation',
  talked_to_users: 'Talked to Users',
  prototype: 'Built Prototype',
  paying_customers: 'Paying Customers',
};

export const STATUS_LABELS: Record<IdeaStatus, string> = {
  buried: 'Buried',
  on_life_support: 'On Life Support',
  resurrected: 'Resurrected',
};

export function formatMoneyCents(cents: number): string {
  const dollars = Math.round(cents / 100);
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(dollars);
}

export function formatShortDate(iso: string): string {
  const datePart = iso.slice(0, 10);
  const [y, m, d] = datePart.split('-').map(Number);
  if (!y || !m || !d) return iso;
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    year: 'numeric',
  }).format(new Date(Date.UTC(y, m - 1, d)));
}

export function formatLifespan(bornIso: string, diedIso: string): string {
  const born = new Date(bornIso.slice(0, 10)).getTime();
  const died = new Date(diedIso.slice(0, 10)).getTime();
  const days = Math.max(0, Math.round((died - born) / (1000 * 60 * 60 * 24)));
  if (days < 7) return `${days}d`;
  if (days < 60) return `${Math.round(days / 7)}w`;
  return `${Math.round(days / 30)}m`;
}
