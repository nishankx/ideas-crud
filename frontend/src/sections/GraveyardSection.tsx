import { AnimatePresence, motion } from 'framer-motion';
import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { apiRequest } from '../api/client';
import { SearchIcon, TombstoneIcon } from '../components/Icons';
import { TombstoneCard } from '../components/TombstoneCard';
import { useAuth } from '../hooks/useAuth';
import { CATEGORY_LABELS, CAUSE_LABELS, STATUS_LABELS } from '../lib/format';
import type { DeathCause, Idea, IdeaCategory, IdeaFilterParams, IdeaStatus, PaginatedIdeas } from '../types/api';

interface GraveyardSectionProps {
  refreshNonce: number;
}

const STATUS_TABS: { value: IdeaStatus | ''; label: string }[] = [
  { value: '', label: 'All Graves' },
  { value: 'buried', label: 'Buried' },
  { value: 'on_life_support', label: 'On Life Support' },
  { value: 'resurrected', label: 'Resurrected' },
];

export function GraveyardSection({ refreshNonce }: GraveyardSectionProps) {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const [query, setQuery] = useState('');
  const [cause, setCause] = useState<DeathCause | ''>('');
  const [category, setCategory] = useState<IdeaCategory | ''>('');
  const [status, setStatus] = useState<IdeaStatus | ''>('');
  const [sort, setSort] = useState<IdeaFilterParams['sort']>('died_at');
  const [page, setPage] = useState(1);

  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [total, setTotal] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [reloadNonce, setReloadNonce] = useState(0);

  useEffect(() => {
    setPage(1);
  }, [query, cause, category, status, sort, isAuthenticated, refreshNonce]);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    const params = new URLSearchParams();
    params.set('page', String(page));
    params.set('limit', '9');
    if (query.trim()) params.set('q', query.trim());
    if (cause) params.set('cause', cause);
    if (category) params.set('category', category);
    if (status) params.set('status', status);
    if (sort) params.set('sort', sort);

    const prefix = isAuthenticated ? '/ideas' : '/demo/ideas';
    apiRequest<PaginatedIdeas>(`${prefix}?${params.toString()}`)
      .then((res) => {
        if (cancelled) return;
        setIdeas((prev) => (page === 1 ? res.items : [...prev, ...res.items]));
        setTotal(res.total);
        setHasMore(res.has_more);
        setLoading(false);
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        setError(err instanceof Error ? err.message : 'Could not load the graveyard');
        setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [page, query, cause, category, status, sort, isAuthenticated, refreshNonce, reloadNonce]);

  return (
    <section id="graveyard" className="mx-auto max-w-6xl px-4 py-14 sm:px-6 sm:py-20">
      <div className="flex flex-col justify-between gap-5 sm:flex-row sm:items-end">
        <div>
          <h2 className="font-display text-3xl font-medium tracking-[-0.02em] text-ink sm:text-4xl">
            The Graveyard
          </h2>
          <p className="mt-1.5 text-sm text-slate-muted">
            {total} {total === 1 ? 'headstone' : 'headstones'} standing in the gallery. Click any stone to inspect its autopsy.
          </p>
        </div>

        {/* Azure Segmented Pill Tabs (matching reference UI) */}
        <div
          role="tablist"
          aria-label="Filter by status"
          className="inline-flex flex-wrap items-center gap-1 rounded-xl border border-hairline bg-alabaster p-1"
        >
          {STATUS_TABS.map((tab) => {
            const active = status === tab.value;
            return (
              <button
                key={tab.label}
                type="button"
                role="tab"
                aria-selected={active}
                onClick={() => setStatus(tab.value)}
                className={`relative rounded-lg px-3 py-1.5 text-xs font-medium transition-colors ${
                  active ? 'text-white' : 'text-slate-muted hover:text-ink'
                }`}
              >
                {active && (
                  <motion.span
                    layoutId="graveyard-status-pill"
                    transition={{ type: 'spring', stiffness: 360, damping: 30 }}
                    className="absolute inset-0 rounded-lg bg-azure shadow-sm"
                  />
                )}
                <span className="relative z-10">{tab.label}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Filter Bar */}
      <div className="mt-6 grid grid-cols-1 gap-3 rounded-2xl border border-hairline bg-alabaster/70 p-3 sm:grid-cols-2 lg:grid-cols-5">
        <label className="relative flex items-center lg:col-span-2">
          <span className="sr-only">Search ideas</span>
          <SearchIcon className="pointer-events-none absolute left-3 h-4 w-4 text-slate-muted" aria-hidden="true" />
          <input
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search epitaphs, names, autopsies..."
            className="w-full rounded-xl border border-hairline bg-canvas py-2 pl-9 pr-3 text-sm text-ink placeholder:text-slate-muted/70 focus:border-azure focus:outline-none"
          />
        </label>

        <select
          aria-label="Filter by cause of death"
          value={cause}
          onChange={(e) => setCause(e.target.value as DeathCause | '')}
          className="rounded-xl border border-hairline bg-canvas px-3 py-2 text-sm text-ink focus:border-azure focus:outline-none"
        >
          <option value="">All Causes</option>
          {Object.entries(CAUSE_LABELS).map(([k, label]) => (
            <option key={k} value={k}>
              {label}
            </option>
          ))}
        </select>

        <select
          aria-label="Filter by category"
          value={category}
          onChange={(e) => setCategory(e.target.value as IdeaCategory | '')}
          className="rounded-xl border border-hairline bg-canvas px-3 py-2 text-sm text-ink focus:border-azure focus:outline-none"
        >
          <option value="">All Categories</option>
          {Object.entries(CATEGORY_LABELS).map(([k, label]) => (
            <option key={k} value={k}>
              {label}
            </option>
          ))}
        </select>

        <select
          aria-label="Sort order"
          value={sort}
          onChange={(e) => setSort(e.target.value as IdeaFilterParams['sort'])}
          className="rounded-xl border border-hairline bg-canvas px-3 py-2 text-sm text-ink focus:border-azure focus:outline-none"
        >
          <option value="died_at">Newest Death</option>
          <option value="money_spent">Most Expensive</option>
          <option value="lifespan">Longest Lived</option>
          <option value="created_at">Recently Buried</option>
        </select>
      </div>

      {/* Error State */}
      {error && (
        <div className="mt-8 rounded-2xl border border-marble-rust/30 bg-marble-rust/5 p-6 text-center">
          <p className="text-sm font-medium text-marble-rust">{error}</p>
          <button
            type="button"
            onClick={() => setReloadNonce((n) => n + 1)}
            className="mt-3 rounded-lg border border-marble-rust/40 bg-canvas px-4 py-1.5 text-xs font-medium text-ink hover:bg-alabaster"
          >
            Try Again
          </button>
        </div>
      )}

      {/* Loading Skeletons */}
      {loading && ideas.length === 0 && !error && (
        <div className="mt-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div
              key={i}
              className="h-64 animate-pulse rounded-t-[2.25rem] rounded-b-2xl border border-hairline bg-alabaster p-6"
            />
          ))}
        </div>
      )}

      {/* Empty State */}
      {!loading && !error && ideas.length === 0 && (
        <div className="mt-10 rounded-3xl border border-dashed border-hairline bg-alabaster/60 p-12 text-center">
          <TombstoneIcon className="mx-auto h-12 w-12 text-slate-muted/50" aria-hidden="true" />
          <h3 className="mt-4 font-display text-xl font-medium text-ink">
            Suspiciously Empty Marble
          </h3>
          <p className="mx-auto mt-2 max-w-md text-sm text-slate-muted">
            No {status ? STATUS_LABELS[status].toLowerCase() + ' ' : ''}ideas match these filters. Either you have a 100% hit rate, or you haven&apos;t confessed yet.
          </p>
          <Link
            to={isAuthenticated ? '/bury' : '/register'}
            className="mt-5 inline-flex rounded-xl bg-azure px-4 py-2 text-xs font-semibold text-white"
          >
            {isAuthenticated ? 'Bury Your First Idea' : 'Create Account to Bury'}
          </Link>
        </div>
      )}

      {/* Tombstone Grid */}
      {ideas.length > 0 && (
        <div className="mt-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          <AnimatePresence mode="popLayout">
            {ideas.map((idea, idx) => (
              <TombstoneCard
                key={idea.id}
                idea={idea}
                index={idx}
                onSelect={(id) => navigate(`/idea/${id}`)}
              />
            ))}
          </AnimatePresence>
        </div>
      )}

      {/* Load More */}
      {hasMore && (
        <div className="mt-10 text-center">
          <button
            type="button"
            disabled={loading}
            onClick={() => setPage((p) => p + 1)}
            className="rounded-xl border border-hairline bg-canvas px-6 py-2.5 text-sm font-medium text-ink shadow-sm transition-colors hover:border-azure/50 hover:bg-alabaster disabled:opacity-50"
          >
            {loading ? 'Exhuming...' : 'Load More Graves'}
          </button>
        </div>
      )}
    </section>
  );
}
