import { AnimatePresence, motion, useReducedMotion } from 'framer-motion';
import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { useFetch } from '../hooks/useFetch';
import { CATEGORY_LABELS, CAUSE_LABELS, VALIDATION_LABELS, formatMoneyCents } from '../lib/format';
import type {
  CausesResponse,
  DeathCause,
  InsightsSummary,
  LifespanResponse,
  PatternsResponse,
  TimelineResponse,
  ValidationResponse,
} from '../types/api';

interface AutopsySectionProps {
  refreshNonce: number;
  summary: InsightsSummary | null;
}

const SMOOTH_EASE = [0.22, 1, 0.36, 1] as const;

type ChartTab = 'causes' | 'timeline' | 'lifespan' | 'validation';

const CHART_TABS: { id: ChartTab; label: string }[] = [
  { id: 'causes', label: 'Causes' },
  { id: 'timeline', label: 'Timeline' },
  { id: 'lifespan', label: 'Lifespan' },
  { id: 'validation', label: 'Validation' },
];

export function AutopsySection({ refreshNonce, summary }: AutopsySectionProps) {
  const { isAuthenticated } = useAuth();
  const shouldReduce = useReducedMotion();
  const [activeTab, setActiveTab] = useState<ChartTab>('causes');

  const base = isAuthenticated ? '/insights' : '/demo/insights';

  const causes = useFetch<CausesResponse>(`${base}/causes?r=${refreshNonce}`);
  const timeline = useFetch<TimelineResponse>(`${base}/timeline?r=${refreshNonce}`);
  const lifespan = useFetch<LifespanResponse>(`${base}/lifespan?r=${refreshNonce}`);
  const validation = useFetch<ValidationResponse>(`${base}/validation?r=${refreshNonce}`);
  const patterns = useFetch<PatternsResponse>(`${base}/patterns?r=${refreshNonce}`);

  const hasError = causes.error || timeline.error || lifespan.error || validation.error || patterns.error;
  const isLoading = causes.loading || timeline.loading || lifespan.loading || validation.loading || patterns.loading;

  const retryAll = () => {
    causes.refetch();
    timeline.refetch();
    lifespan.refetch();
    validation.refetch();
    patterns.refetch();
  };

  return (
    <section
      id="autopsy"
      className="mx-auto max-w-6xl border-t border-hairline px-4 py-16 sm:px-6 sm:py-24"
    >
      {/* Azure Tab Strip + Split Feature Header (inspired by reference) */}
      <div className="flex flex-col justify-between gap-6 lg:flex-row lg:items-end">
        <div className="max-w-xl">
          <div className="inline-flex rounded-xl bg-azure p-1 shadow-sm">
            {CHART_TABS.map((tab) => {
              const active = activeTab === tab.id;
              return (
                <button
                  key={tab.id}
                  type="button"
                  onClick={() => setActiveTab(tab.id)}
                  className={`relative rounded-lg px-3.5 py-1.5 text-xs font-semibold transition-colors ${
                    active ? 'text-azure-deep' : 'text-white/85 hover:text-white'
                  }`}
                >
                  {active && (
                    <motion.span
                      layoutId="autopsy-tab-pill"
                      transition={{ type: 'spring', stiffness: 360, damping: 30 }}
                      className="absolute inset-0 rounded-lg bg-white shadow-sm"
                    />
                  )}
                  <span className="relative z-10">{tab.label}</span>
                </button>
              );
            })}
          </div>

          <h2 className="mt-5 font-display text-3xl font-medium tracking-[-0.02em] text-ink sm:text-4xl">
            Autopsy Report
          </h2>
          <p className="mt-2 text-sm leading-relaxed text-slate-muted">
            Forensic patterns across your graveyard — why ideas keep dying, how fast they flatline, and whether talking to users saved them.
          </p>
        </div>

        {/* Headline Stats Grid */}
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
          <StatCard label="Total Buried" value={summary ? String(summary.total_ideas) : '—'} />
          <StatCard
            label="Money Lost"
            value={summary ? formatMoneyCents(summary.total_money_lost_cents) : '—'}
          />
          <StatCard label="Hours Lost" value={summary ? `${summary.total_hours_lost}h` : '—'} />
          <StatCard
            label="Avg Lifespan"
            value={summary ? `${Math.round(summary.avg_lifespan_days)}d` : '—'}
          />
          <StatCard
            label="Top Killer"
            value={summary?.most_common_cause ? CAUSE_LABELS[summary.most_common_cause] : '—'}
            accent="rust"
          />
          <StatCard
            label="Deadliest Niche"
            value={summary?.most_lethal_category ? CATEGORY_LABELS[summary.most_lethal_category] : '—'}
          />
        </div>
      </div>

      {hasError && (
        <div className="mt-8 rounded-2xl border border-marble-rust/30 bg-marble-rust/5 p-6 text-center">
          <p className="text-sm font-medium text-marble-rust">Could not load autopsy analytics.</p>
          <button
            type="button"
            onClick={retryAll}
            className="mt-3 rounded-lg border border-marble-rust/40 bg-canvas px-4 py-1.5 text-xs font-medium text-ink hover:bg-alabaster"
          >
            Retry Report
          </button>
        </div>
      )}

      {isLoading && !hasError && (
        <div className="mt-10 grid grid-cols-1 gap-6 lg:grid-cols-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-64 animate-pulse rounded-2xl border border-hairline bg-alabaster" />
          ))}
        </div>
      )}

      {!isLoading && !hasError && (
        <div className="mt-10 space-y-10">
          {/* Interactive Tab Spotlight Card */}
          <div className="rounded-3xl border border-hairline bg-alabaster/70 p-6 sm:p-8">
            <AnimatePresence mode="wait">
              {activeTab === 'causes' && (
                <motion.div
                  key="causes"
                  initial={shouldReduce ? false : { opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={shouldReduce ? undefined : { opacity: 0, y: -8 }}
                  transition={{ duration: 0.28, ease: SMOOTH_EASE }}
                >
                  <CausesPanel causes={causes.data} />
                </motion.div>
              )}
              {activeTab === 'timeline' && (
                <motion.div
                  key="timeline"
                  initial={shouldReduce ? false : { opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={shouldReduce ? undefined : { opacity: 0, y: -8 }}
                  transition={{ duration: 0.28, ease: SMOOTH_EASE }}
                >
                  <h3 className="font-display text-xl font-semibold text-ink">Monthly Burial Timeline</h3>
                  <p className="mt-1 text-xs text-slate-muted">Ideas laid to rest per month</p>
                  <div className="mt-6">
                    <TimelineSVG buckets={timeline.data?.buckets ?? []} />
                  </div>
                </motion.div>
              )}
              {activeTab === 'lifespan' && (
                <motion.div
                  key="lifespan"
                  initial={shouldReduce ? false : { opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={shouldReduce ? undefined : { opacity: 0, y: -8 }}
                  transition={{ duration: 0.28, ease: SMOOTH_EASE }}
                >
                  <LifespanPanel lifespan={lifespan.data} />
                </motion.div>
              )}
              {activeTab === 'validation' && (
                <motion.div
                  key="validation"
                  initial={shouldReduce ? false : { opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={shouldReduce ? undefined : { opacity: 0, y: -8 }}
                  transition={{ duration: 0.28, ease: SMOOTH_EASE }}
                >
                  <h3 className="font-display text-xl font-semibold text-ink">Validation vs. Cause Matrix</h3>
                  <p className="mt-1 text-xs text-slate-muted">
                    Does talking to users actually prevent No Market Need?
                  </p>
                  <ValidationHeatmap
                    matrix={validation.data?.matrix ?? []}
                    topCauses={(causes.data?.causes ?? []).slice(0, 4).map((c) => c.cause)}
                  />
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {/* Patterns We Noticed */}
          {patterns.data && patterns.data.patterns.length > 0 && (
            <div>
              <h3 className="font-display text-2xl font-medium text-ink">
                Patterns We Noticed
              </h3>
              <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                {patterns.data.patterns.map((p, idx) => {
                  const accent =
                    p.severity === 'critical'
                      ? 'border-marble-rust/35 bg-marble-rust/5'
                      : p.severity === 'warning'
                        ? 'border-amber/35 bg-amber/5'
                        : 'border-azure/30 bg-cloud-wash/70';
                  return (
                    <motion.div
                      key={p.id}
                      initial={shouldReduce ? false : { opacity: 0, y: 14 }}
                      whileInView={{ opacity: 1, y: 0 }}
                      viewport={{ once: true, amount: 0.2 }}
                      transition={{ duration: 0.45, delay: idx * 0.06, ease: SMOOTH_EASE }}
                      className={`rounded-2xl border p-5 ${accent}`}
                    >
                      <div className="flex items-center justify-between gap-2">
                        <span className="font-display text-base font-semibold text-ink">
                          {p.title}
                        </span>
                        <span className="rounded-md bg-canvas px-2 py-0.5 font-mono text-[11px] text-slate-muted">
                          n={p.sample_size}
                        </span>
                      </div>
                      <p className="mt-2 text-sm leading-relaxed text-ink/85">{p.message}</p>
                    </motion.div>
                  );
                })}
              </div>
            </div>
          )}
        </div>
      )}
    </section>
  );
}

function CausesPanel({ causes }: { causes: CausesResponse | null }) {
  const shouldReduce = useReducedMotion();
  return (
    <div>
      <h3 className="font-display text-xl font-semibold text-ink">Causes of Death</h3>
      <p className="mt-1 text-xs text-slate-muted">Primary killers ranked by burial count</p>

      <div className="mt-6 space-y-4">
        {(causes?.causes ?? []).length === 0 ? (
          <p className="text-xs text-slate-muted">No burials recorded yet.</p>
        ) : (
          causes?.causes.map((item) => {
            const pct = Math.round(item.share * 100);
            return (
              <div key={item.cause}>
                <div className="flex items-center justify-between text-xs">
                  <span className="font-medium text-ink">{CAUSE_LABELS[item.cause]}</span>
                  <span className="font-mono text-slate-muted">
                    {item.count} ({pct}%)
                  </span>
                </div>
                <div className="mt-1.5 h-2.5 w-full overflow-hidden rounded-full bg-hairline">
                  <motion.div
                    initial={shouldReduce ? { width: `${pct}%` } : { width: 0 }}
                    animate={{ width: `${pct}%` }}
                    transition={{ duration: 0.65, ease: SMOOTH_EASE }}
                    className="h-full rounded-full bg-azure"
                  />
                </div>
                {item.co_occurring.length > 0 && (
                  <p className="mt-1 text-[11px] text-slate-muted">
                    Often paired with:{' '}
                    {item.co_occurring
                      .slice(0, 2)
                      .map((c) => `${CAUSE_LABELS[c.cause]} (${c.count})`)
                      .join(', ')}
                  </p>
                )}
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}

function LifespanPanel({ lifespan }: { lifespan: LifespanResponse | null }) {
  const shouldReduce = useReducedMotion();
  return (
    <div>
      <h3 className="font-display text-xl font-semibold text-ink">Lifespan Distribution</h3>
      <p className="mt-1 text-xs text-slate-muted">How long ideas survive before flatlining</p>

      <div className="mt-6 grid h-48 grid-cols-5 items-end gap-3 sm:gap-5">
        {(lifespan?.buckets ?? []).map((b) => {
          const pct = Math.max(8, Math.round(b.share * 100));
          return (
            <div key={b.bucket} className="flex h-full flex-col items-center justify-end gap-2">
              <span className="font-mono text-xs font-semibold text-ink">{b.count}</span>
              <div className="flex w-full flex-1 items-end justify-center rounded-xl border border-hairline bg-canvas p-1.5">
                <motion.div
                  initial={shouldReduce ? { height: `${pct}%` } : { height: 0 }}
                  animate={{ height: `${pct}%` }}
                  transition={{ duration: 0.6, ease: SMOOTH_EASE }}
                  className="w-full rounded-lg bg-gradient-to-t from-azure-deep to-[#6BB5FA]"
                />
              </div>
              <span className="text-center font-mono text-[11px] text-slate-muted">
                {b.bucket}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
  accent,
}: {
  label: string;
  value: string;
  accent?: 'rust';
}) {
  const color = accent === 'rust' ? 'text-marble-rust' : 'text-ink';
  return (
    <div className="rounded-xl border border-hairline bg-canvas p-3.5 shadow-sm">
      <div className={`truncate font-mono text-sm font-semibold sm:text-base ${color}`} title={value}>
        {value}
      </div>
      <div className="mt-0.5 text-[10px] font-medium uppercase tracking-wider text-slate-muted">
        {label}
      </div>
    </div>
  );
}

function TimelineSVG({ buckets }: { buckets: { month: string; count: number }[] }) {
  const shouldReduce = useReducedMotion();

  if (buckets.length === 0) {
    return <p className="py-12 text-center text-xs text-slate-muted">No timeline data yet.</p>;
  }

  const width = 560;
  const height = 190;
  const padX = 24;
  const padY = 24;
  const maxCount = Math.max(1, ...buckets.map((b) => b.count));

  const points = buckets.map((b, i) => {
    const x =
      buckets.length === 1
        ? width / 2
        : padX + (i / (buckets.length - 1)) * (width - padX * 2);
    const y = height - padY - (b.count / maxCount) * (height - padY * 2);
    return { x, y, ...b };
  });

  const pathD = points
    .map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(1)} ${p.y.toFixed(1)}`)
    .join(' ');

  return (
    <div className="rounded-2xl border border-hairline bg-canvas p-4">
      <svg
        viewBox={`0 0 ${width} ${height}`}
        role="img"
        aria-label="Monthly burial timeline chart"
        className="h-44 w-full overflow-visible"
      >
        <line x1={padX} y1={height - padY} x2={width - padX} y2={height - padY} stroke="#E4E7EC" strokeWidth="1" />
        <motion.path
          d={pathD}
          pathLength={1}
          strokeDasharray="1"
          fill="none"
          stroke="#2589F5"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          initial={shouldReduce ? false : { strokeDashoffset: 1 }}
          animate={{ strokeDashoffset: 0 }}
          transition={{ duration: 0.9, ease: SMOOTH_EASE }}
        />
        {points.map((p) => (
          <g key={p.month}>
            <circle cx={p.x} cy={p.y} r="4.5" fill="#FFFFFF" stroke="#2589F5" strokeWidth="2.5" />
            <title>{`${p.month}: ${p.count} buried`}</title>
          </g>
        ))}
      </svg>
      <div className="mt-2 flex justify-between font-mono text-[10px] text-slate-muted">
        <span>{buckets[0]?.month}</span>
        {buckets.length > 2 && <span>{buckets[Math.floor(buckets.length / 2)]?.month}</span>}
        <span>{buckets[buckets.length - 1]?.month}</span>
      </div>
    </div>
  );
}

function ValidationHeatmap({
  matrix,
  topCauses,
}: {
  matrix: ValidationResponse['matrix'];
  topCauses: DeathCause[];
}) {
  if (matrix.length === 0 || topCauses.length === 0) {
    return <p className="py-12 text-center text-xs text-slate-muted">Not enough data for matrix yet.</p>;
  }

  let maxCell = 1;
  for (const row of matrix) {
    for (const c of topCauses) {
      maxCell = Math.max(maxCell, row.causes[c] ?? 0);
    }
  }

  return (
    <div className="mt-5 overflow-x-auto rounded-2xl border border-hairline bg-canvas p-4">
      <table className="w-full border-collapse text-left text-xs">
        <thead>
          <tr>
            <th className="pb-2 pr-2 font-medium text-slate-muted">Validation</th>
            {topCauses.map((c) => (
              <th key={c} className="px-1 pb-2 text-center font-mono text-[10px] text-slate-muted">
                {CAUSE_LABELS[c].split(' ').slice(0, 2).join(' ')}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-hairline">
          {matrix.map((row) => (
            <tr key={row.validation_level}>
              <td className="py-2.5 pr-2 font-medium text-ink">
                {VALIDATION_LABELS[row.validation_level]}
              </td>
              {topCauses.map((c) => {
                const val = row.causes[c] ?? 0;
                const alpha = val === 0 ? 0.04 : 0.14 + (val / maxCell) * 0.65;
                return (
                  <td key={c} className="p-1 text-center">
                    <div
                      className="mx-auto flex h-8 w-11 items-center justify-center rounded-lg border border-hairline font-mono text-xs font-semibold text-ink"
                      style={{ backgroundColor: `rgba(37, 137, 245, ${alpha.toFixed(2)})` }}
                      title={`${VALIDATION_LABELS[row.validation_level]} × ${CAUSE_LABELS[c]}: ${val}`}
                    >
                      {val}
                    </div>
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
