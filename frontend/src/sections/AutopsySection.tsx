import { motion, useReducedMotion } from 'framer-motion';
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

export function AutopsySection({ refreshNonce, summary }: AutopsySectionProps) {
  const { isAuthenticated } = useAuth();
  const shouldReduce = useReducedMotion();
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
      className="mx-auto max-w-6xl border-t border-stone-border/60 px-4 py-16 sm:px-6 sm:py-24"
    >
      <div className="max-w-2xl">
        <h2 className="font-display text-3xl font-bold text-bone sm:text-4xl">
          Autopsy Report
        </h2>
        <p className="mt-2 text-sm text-stone">
          Forensic patterns across your graveyard — why ideas keep dying, how fast they flatline, and whether validation saved them.
        </p>
      </div>

      {/* Headline Stats Row */}
      <div className="mt-8 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
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

      {hasError && (
        <div className="mt-8 rounded-xl border border-rust/40 bg-rust/10 p-6 text-center">
          <p className="text-sm text-rust">Could not load autopsy analytics.</p>
          <button
            type="button"
            onClick={retryAll}
            className="mt-3 rounded-lg border border-rust/50 px-4 py-1.5 text-xs text-bone hover:bg-rust/20"
          >
            Retry Report
          </button>
        </div>
      )}

      {isLoading && !hasError && (
        <div className="mt-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-64 animate-pulse rounded-2xl border border-stone-border bg-headstone/70" />
          ))}
        </div>
      )}

      {!isLoading && !hasError && (
        <div className="mt-10 space-y-8">
          {/* Patterns We Noticed */}
          {patterns.data && patterns.data.patterns.length > 0 && (
            <div>
              <h3 className="font-display text-xl font-semibold text-bone">
                Patterns We Noticed
              </h3>
              <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                {patterns.data.patterns.map((p, idx) => {
                  const border =
                    p.severity === 'critical'
                      ? 'border-rust/50 bg-rust/5'
                      : p.severity === 'warning'
                        ? 'border-amber/40 bg-amber/5'
                        : 'border-necro/40 bg-necro/5';
                  return (
                    <motion.div
                      key={p.id}
                      initial={shouldReduce ? false : { opacity: 0, y: 16 }}
                      whileInView={{ opacity: 1, y: 0 }}
                      viewport={{ once: true, amount: 0.2 }}
                      transition={{ duration: 0.35, delay: idx * 0.06 }}
                      className={`rounded-xl border p-5 ${border}`}
                    >
                      <div className="flex items-center justify-between gap-2">
                        <span className="font-display text-sm font-semibold text-bone">
                          {p.title}
                        </span>
                        <span className="font-mono text-[11px] text-stone">
                          n={p.sample_size}
                        </span>
                      </div>
                      <p className="mt-2 text-sm leading-relaxed text-bone/90">{p.message}</p>
                    </motion.div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Charts Grid */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {/* 1. Cause of Death Breakdown */}
            <div className="rounded-2xl border border-stone-border bg-headstone p-6">
              <h3 className="font-display text-lg font-semibold text-bone">
                Causes of Death
              </h3>
              <p className="mt-1 text-xs text-stone">
                Primary killers ranked by burial count
              </p>

              <div className="mt-5 space-y-3.5">
                {(causes.data?.causes ?? []).length === 0 ? (
                  <p className="text-xs text-stone">No burials recorded yet.</p>
                ) : (
                  causes.data?.causes.map((item) => {
                    const pct = Math.round(item.share * 100);
                    return (
                      <div key={item.cause}>
                        <div className="flex items-center justify-between text-xs">
                          <span className="font-medium text-bone">
                            {CAUSE_LABELS[item.cause]}
                          </span>
                          <span className="font-mono text-stone">
                            {item.count} ({pct}%)
                          </span>
                        </div>
                        <div className="mt-1.5 h-2.5 w-full overflow-hidden rounded-full bg-crypt">
                          <motion.div
                            initial={shouldReduce ? { width: `${pct}%` } : { width: 0 }}
                            whileInView={{ width: `${pct}%` }}
                            viewport={{ once: true, amount: 0.3 }}
                            transition={{ duration: 0.6, ease: 'easeOut' }}
                            className="h-full rounded-full bg-necro"
                          />
                        </div>
                        {item.co_occurring.length > 0 && (
                          <p className="mt-1 text-[11px] text-stone/80">
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

            {/* 2. Monthly Burial Timeline (SVG) */}
            <div className="rounded-2xl border border-stone-border bg-headstone p-6">
              <h3 className="font-display text-lg font-semibold text-bone">
                Monthly Burial Timeline
              </h3>
              <p className="mt-1 text-xs text-stone">
                Ideas laid to rest per month
              </p>

              <div className="mt-5">
                <TimelineSVG buckets={timeline.data?.buckets ?? []} />
              </div>
            </div>

            {/* 3. Lifespan Distribution */}
            <div className="rounded-2xl border border-stone-border bg-headstone p-6">
              <h3 className="font-display text-lg font-semibold text-bone">
                Lifespan Distribution
              </h3>
              <p className="mt-1 text-xs text-stone">
                How long ideas survive before flatlining
              </p>

              <div className="mt-6 grid grid-cols-5 items-end gap-3 h-44">
                {(lifespan.data?.buckets ?? []).map((b) => {
                  const pct = Math.max(6, Math.round(b.share * 100));
                  return (
                    <div key={b.bucket} className="flex h-full flex-col items-center justify-end gap-2">
                      <span className="font-mono text-xs text-bone">{b.count}</span>
                      <div className="flex w-full flex-1 items-end justify-center rounded-t-lg bg-crypt/70 p-1">
                        <motion.div
                          initial={shouldReduce ? { height: `${pct}%` } : { height: 0 }}
                          whileInView={{ height: `${pct}%` }}
                          viewport={{ once: true, amount: 0.2 }}
                          transition={{ duration: 0.55, ease: 'easeOut' }}
                          className="w-full rounded-t-md bg-gradient-to-t from-necro-dim to-necro"
                        />
                      </div>
                      <span className="text-center font-mono text-[11px] text-stone">
                        {b.bucket}
                      </span>
                    </div>
                  );
                })}
              </div>
            </div>

            {/* 4. Validation vs Cause Heatmap Matrix */}
            <div className="rounded-2xl border border-stone-border bg-headstone p-6">
              <h3 className="font-display text-lg font-semibold text-bone">
                Validation vs. Cause Matrix
              </h3>
              <p className="mt-1 text-xs text-stone">
                Does talking to users actually prevent No Market Need?
              </p>

              <ValidationHeatmap
                matrix={validation.data?.matrix ?? []}
                topCauses={(causes.data?.causes ?? []).slice(0, 4).map((c) => c.cause)}
              />
            </div>
          </div>
        </div>
      )}
    </section>
  );
}

function StatCard({
  label,
  value,
  accent,
}: {
  label: string;
  value: string;
  accent?: 'rust' | 'necro';
}) {
  const color = accent === 'rust' ? 'text-rust' : accent === 'necro' ? 'text-necro' : 'text-bone';
  return (
    <div className="rounded-xl border border-stone-border bg-headstone p-4">
      <div className={`truncate font-mono text-base font-semibold sm:text-lg ${color}`} title={value}>
        {value}
      </div>
      <div className="mt-1 text-[11px] uppercase tracking-wider text-stone">{label}</div>
    </div>
  );
}

function TimelineSVG({ buckets }: { buckets: { month: string; count: number }[] }) {
  const shouldReduce = useReducedMotion();

  if (buckets.length === 0) {
    return <p className="py-12 text-center text-xs text-stone">No timeline data yet.</p>;
  }

  const width = 480;
  const height = 180;
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
    <div>
      <svg
        viewBox={`0 0 ${width} ${height}`}
        role="img"
        aria-label="Monthly burial timeline chart"
        className="h-44 w-full overflow-visible"
      >
        <line x1={padX} y1={height - padY} x2={width - padX} y2={height - padY} stroke="#272C37" strokeWidth="1" />
        <motion.path
          d={pathD}
          pathLength={1}
          strokeDasharray="1"
          fill="none"
          stroke="#7CEB4F"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          initial={shouldReduce ? false : { strokeDashoffset: 1 }}
          whileInView={{ strokeDashoffset: 0 }}
          viewport={{ once: true, amount: 0.3 }}
          transition={{ duration: 0.9, ease: 'easeInOut' }}
        />
        {points.map((p) => (
          <g key={p.month}>
            <circle cx={p.x} cy={p.y} r="4" fill="#0A0B0D" stroke="#7CEB4F" strokeWidth="2" />
            <title>{`${p.month}: ${p.count} buried`}</title>
          </g>
        ))}
      </svg>
      <div className="mt-2 flex justify-between font-mono text-[10px] text-stone">
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
    return <p className="py-12 text-center text-xs text-stone">Not enough data for matrix yet.</p>;
  }

  let maxCell = 1;
  for (const row of matrix) {
    for (const c of topCauses) {
      maxCell = Math.max(maxCell, row.causes[c] ?? 0);
    }
  }

  return (
    <div className="mt-5 overflow-x-auto">
      <table className="w-full border-collapse text-left text-xs">
        <thead>
          <tr>
            <th className="pb-2 pr-2 font-medium text-stone">Validation</th>
            {topCauses.map((c) => (
              <th key={c} className="pb-2 px-1 text-center font-mono text-[10px] text-stone">
                {CAUSE_LABELS[c].split(' ').slice(0, 2).join(' ')}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-stone-border/40">
          {matrix.map((row) => (
            <tr key={row.validation_level}>
              <td className="py-2 pr-2 font-medium text-bone">
                {VALIDATION_LABELS[row.validation_level]}
              </td>
              {topCauses.map((c) => {
                const val = row.causes[c] ?? 0;
                const alpha = val === 0 ? 0.04 : 0.15 + (val / maxCell) * 0.65;
                return (
                  <td key={c} className="p-1 text-center">
                    <div
                      className="mx-auto flex h-8 w-10 items-center justify-center rounded-md border border-stone-border/50 font-mono text-xs text-bone"
                      style={{ backgroundColor: `rgba(124, 235, 79, ${alpha.toFixed(2)})` }}
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
