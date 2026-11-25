import { motion, useReducedMotion } from 'framer-motion';
import { BoltIcon, SkullIcon, TombstoneIcon } from './Icons';
import { CATEGORY_LABELS, CAUSE_LABELS, STATUS_LABELS, formatLifespan, formatShortDate } from '../lib/format';
import type { Idea } from '../types/api';

interface TombstoneCardProps {
  idea: Idea;
  index: number;
  onSelect: (id: string) => void;
}

const STATUS_STYLE: Record<Idea['status'], string> = {
  buried: 'border-hairline bg-alabaster text-slate-muted',
  on_life_support: 'border-amber/30 bg-amber/10 text-amber',
  resurrected: 'border-revive/30 bg-revive/10 text-revive',
};

const SMOOTH_EASE = [0.22, 1, 0.36, 1] as const;

export function TombstoneCard({ idea, index, onSelect }: TombstoneCardProps) {
  const shouldReduce = useReducedMotion();

  return (
    <motion.article
      layout="position"
      layoutId={`tombstone-${idea.id}`}
      initial={shouldReduce ? false : { opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount: 0.12 }}
      transition={{
        duration: 0.5,
        delay: Math.min(index * 0.04, 0.2),
        ease: SMOOTH_EASE,
      }}
      whileHover={shouldReduce ? undefined : { y: -4, transition: { duration: 0.22, ease: SMOOTH_EASE } }}
      className="group relative flex flex-col justify-between rounded-t-[2.25rem] rounded-b-2xl border border-hairline bg-canvas p-6 shadow-tombstone transition-all duration-300 hover:border-azure/50 hover:shadow-azure-glow"
    >
      <button
        type="button"
        onClick={() => onSelect(idea.id)}
        className="text-left focus:outline-none"
        aria-label={`Inspect tombstone for ${idea.name}`}
      >
        <div className="flex items-center justify-between gap-2">
          <span className="rounded-md border border-hairline bg-alabaster px-2 py-0.5 font-mono text-[11px] uppercase tracking-wider text-slate-muted">
            {CATEGORY_LABELS[idea.category]}
          </span>

          <span
            className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 font-mono text-[11px] font-medium ${STATUS_STYLE[idea.status]}`}
          >
            {idea.status === 'resurrected' && <BoltIcon className="h-3 w-3" aria-hidden="true" />}
            {STATUS_LABELS[idea.status]}
          </span>
        </div>

        <div className="mt-5 text-center">
          <TombstoneIcon
            className="mx-auto mb-2 h-6 w-6 text-slate-muted/50 transition-colors duration-300 group-hover:text-azure"
            aria-hidden="true"
          />
          <h3 className="font-display text-xl font-semibold text-ink transition-colors duration-300 group-hover:text-azure-deep">
            {idea.name}
          </h3>
          <p className="mt-1 font-mono text-xs text-slate-muted">
            {formatShortDate(idea.born_at)} — {formatShortDate(idea.died_at)} ({formatLifespan(idea.born_at, idea.died_at)})
          </p>
        </div>

        <blockquote className="mt-4 border-y border-hairline py-3 text-center font-display text-sm italic leading-relaxed text-ink/85">
          “{idea.epitaph}”
        </blockquote>
      </button>

      <div className="mt-4 flex items-center justify-between text-xs text-slate-muted">
        <span className="inline-flex items-center gap-1.5 font-medium text-marble-rust">
          <SkullIcon className="h-3.5 w-3.5" aria-hidden="true" />
          <span>{CAUSE_LABELS[idea.primary_cause]}</span>
        </span>
        <span className="font-mono text-[11px] text-slate-muted">
          {idea.hours_invested}h lost
        </span>
      </div>
    </motion.article>
  );
}
