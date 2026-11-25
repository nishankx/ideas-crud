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
  buried: 'border-stone-border bg-slate-mist/70 text-stone',
  on_life_support: 'border-amber/40 bg-amber/10 text-amber',
  resurrected: 'border-necro/50 bg-necro/15 text-necro',
};

export function TombstoneCard({ idea, index, onSelect }: TombstoneCardProps) {
  const shouldReduce = useReducedMotion();

  return (
    <motion.article
      layout
      layoutId={`tombstone-${idea.id}`}
      initial={shouldReduce ? false : { opacity: 0, y: 28 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount: 0.15 }}
      transition={{ duration: 0.4, delay: Math.min(index * 0.04, 0.25) }}
      whileHover={shouldReduce ? undefined : { y: -4 }}
      className="group relative flex flex-col justify-between rounded-t-[2.25rem] rounded-b-xl border border-stone-border bg-headstone p-6 shadow-tombstone transition-colors hover:border-necro/60 hover:shadow-necro-glow"
    >
      <button
        type="button"
        onClick={() => onSelect(idea.id)}
        className="text-left focus:outline-none"
        aria-label={`Inspect tombstone for ${idea.name}`}
      >
        <div className="flex items-center justify-between gap-2">
          <span className="rounded-md border border-stone-border bg-crypt px-2 py-0.5 font-mono text-[11px] uppercase tracking-wider text-stone">
            {CATEGORY_LABELS[idea.category]}
          </span>

          <span
            className={`inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 font-mono text-[11px] ${STATUS_STYLE[idea.status]}`}
          >
            {idea.status === 'resurrected' && <BoltIcon className="h-3 w-3" aria-hidden="true" />}
            {STATUS_LABELS[idea.status]}
          </span>
        </div>

        <div className="mt-5 text-center">
          <TombstoneIcon
            className="mx-auto mb-2 h-6 w-6 text-stone/60 transition-colors group-hover:text-necro"
            aria-hidden="true"
          />
          <h3 className="font-display text-xl font-bold text-bone group-hover:text-necro transition-colors">
            {idea.name}
          </h3>
          <p className="mt-1 font-mono text-xs text-stone">
            {formatShortDate(idea.born_at)} — {formatShortDate(idea.died_at)} ({formatLifespan(idea.born_at, idea.died_at)})
          </p>
        </div>

        <blockquote className="mt-4 border-y border-stone-border/60 py-3 text-center font-display text-sm italic leading-relaxed text-bone/90">
          “{idea.epitaph}”
        </blockquote>
      </button>

      <div className="mt-4 flex items-center justify-between text-xs text-stone">
        <span className="inline-flex items-center gap-1.5 text-rust">
          <SkullIcon className="h-3.5 w-3.5" aria-hidden="true" />
          <span>{CAUSE_LABELS[idea.primary_cause]}</span>
        </span>
        <span className="font-mono text-[11px] text-stone/80">
          {idea.hours_invested}h lost
        </span>
      </div>
    </motion.article>
  );
}
