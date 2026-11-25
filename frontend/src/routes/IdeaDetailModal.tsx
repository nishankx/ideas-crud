import { motion, useReducedMotion } from 'framer-motion';
import { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ApiError, apiRequest } from '../api/client';
import { BoltIcon, SkullIcon } from '../components/Icons';
import { Modal } from '../components/Modal';
import { useAuth } from '../hooks/useAuth';
import { useFetch } from '../hooks/useFetch';
import {
  CATEGORY_LABELS,
  CAUSE_LABELS,
  STATUS_LABELS,
  VALIDATION_LABELS,
  formatLifespan,
  formatMoneyCents,
  formatShortDate,
} from '../lib/format';
import type { Idea } from '../types/api';

interface IdeaDetailModalProps {
  onMutated: () => void;
}

export function IdeaDetailModal({ onMutated }: IdeaDetailModalProps) {
  const { id } = useParams<{ id: string }>();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const shouldReduce = useReducedMotion();

  const endpoint = id ? (isAuthenticated ? `/ideas/${id}` : `/demo/ideas/${id}`) : null;
  const { data: idea, loading, error, refetch, setData } = useFetch<Idea>(endpoint);

  const [actionError, setActionError] = useState<string | null>(null);
  const [cracking, setCracking] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);

  const close = () => navigate('/');

  const handleResurrect = async () => {
    if (!idea || !isAuthenticated) return;
    setActionError(null);
    setCracking(true);

    const nextStatus = idea.status === 'buried' ? 'on_life_support' : 'resurrected';
    const prev = idea;
    setData({ ...idea, status: nextStatus });

    try {
      const updated = await apiRequest<Idea>(`/ideas/${idea.id}/resurrect`, {
        method: 'POST',
      });
      setData(updated);
      onMutated();
    } catch (err) {
      setData(prev);
      setActionError(err instanceof ApiError ? err.message : 'Resurrection failed.');
    } finally {
      setTimeout(() => setCracking(false), shouldReduce ? 0 : 600);
    }
  };

  const handleDelete = async () => {
    if (!idea || !isAuthenticated) return;
    try {
      await apiRequest<void>(`/ideas/${idea.id}`, { method: 'DELETE' });
      onMutated();
      navigate('/');
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : 'Failed to exhume & delete.');
    }
  };

  return (
    <Modal
      title={idea?.name ?? 'Tombstone Inspection'}
      titleId="idea-detail-title"
      onClose={close}
      maxWidthClass="max-w-xl"
    >
      {loading && (
        <div className="space-y-3 py-6">
          <div className="h-6 w-2/3 animate-pulse rounded-lg bg-alabaster" />
          <div className="h-20 animate-pulse rounded-lg bg-alabaster" />
          <div className="h-24 animate-pulse rounded-lg bg-alabaster" />
        </div>
      )}

      {error && (
        <div className="py-6 text-center">
          <p className="text-sm text-marble-rust">{error}</p>
          <button
            type="button"
            onClick={refetch}
            className="mt-3 rounded-xl border border-hairline px-4 py-1.5 text-xs font-medium text-ink"
          >
            Retry
          </button>
        </div>
      )}

      {idea && (
        <motion.div
          layoutId={`tombstone-${idea.id}`}
          animate={
            cracking && !shouldReduce
              ? {
                  boxShadow: [
                    '0 0 0px rgba(37,137,245,0)',
                    '0 0 36px rgba(37,137,245,0.45)',
                    '0 0 0px rgba(37,137,245,0)',
                  ],
                  x: [0, -3, 3, -2, 2, 0],
                }
              : undefined
          }
          transition={{ duration: 0.55, ease: [0.22, 1, 0.36, 1] }}
          className="space-y-5"
        >
          {actionError && (
            <div role="alert" className="rounded-xl border border-marble-rust/30 bg-marble-rust/5 px-3 py-2 text-xs text-marble-rust">
              {actionError}
            </div>
          )}

          <div className="flex flex-wrap items-center gap-2 text-xs">
            <span className="rounded-md border border-hairline bg-alabaster px-2.5 py-1 font-mono uppercase text-slate-muted">
              {CATEGORY_LABELS[idea.category]}
            </span>
            <span className="rounded-full border border-azure/30 bg-cloud-wash px-2.5 py-1 font-mono font-medium text-azure-deep">
              {STATUS_LABELS[idea.status]}
            </span>
            <span className="font-mono text-slate-muted">
              {formatShortDate(idea.born_at)} — {formatShortDate(idea.died_at)} ({formatLifespan(idea.born_at, idea.died_at)})
            </span>
          </div>

          {idea.tagline && <p className="text-sm text-slate-muted">{idea.tagline}</p>}

          <blockquote className="rounded-2xl border border-hairline bg-alabaster p-4 text-center font-display text-base italic text-ink">
            “{idea.epitaph}”
          </blockquote>

          <div className="grid grid-cols-3 gap-3 text-center">
            <div className="rounded-xl border border-hairline bg-alabaster p-3">
              <div className="font-mono text-sm font-semibold text-ink">
                {formatMoneyCents(idea.money_spent_cents)}
              </div>
              <div className="text-[11px] text-slate-muted">Money Lost</div>
            </div>
            <div className="rounded-xl border border-hairline bg-alabaster p-3">
              <div className="font-mono text-sm font-semibold text-ink">
                {idea.hours_invested}h
              </div>
              <div className="text-[11px] text-slate-muted">Hours Invested</div>
            </div>
            <div className="rounded-xl border border-hairline bg-alabaster p-3">
              <div className="font-mono text-xs font-semibold text-ink">
                {VALIDATION_LABELS[idea.validation_level]}
              </div>
              <div className="text-[11px] text-slate-muted">Validation</div>
            </div>
          </div>

          <div>
            <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-muted">
              Cause of Death
            </h3>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <span className="inline-flex items-center gap-1.5 rounded-full border border-marble-rust/30 bg-marble-rust/5 px-3 py-1 text-xs font-medium text-marble-rust">
                <SkullIcon className="h-3.5 w-3.5" aria-hidden="true" />
                {CAUSE_LABELS[idea.primary_cause]}
              </span>
              {idea.secondary_causes.map((sc) => (
                <span
                  key={sc}
                  className="rounded-full border border-hairline bg-alabaster px-2.5 py-1 text-xs text-slate-muted"
                >
                  + {CAUSE_LABELS[sc]}
                </span>
              ))}
            </div>
          </div>

          {idea.autopsy && (
            <div>
              <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-muted">
                Autopsy Notes
              </h3>
              <p className="mt-1.5 whitespace-pre-line text-sm leading-relaxed text-ink/90">
                {idea.autopsy}
              </p>
            </div>
          )}

          {idea.description && (
            <div>
              <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-muted">
                What It Was
              </h3>
              <p className="mt-1.5 whitespace-pre-line text-sm text-slate-muted">
                {idea.description}
              </p>
            </div>
          )}

          {isAuthenticated && (
            <div className="flex flex-wrap items-center justify-between gap-3 border-t border-hairline pt-4">
              <div>
                {confirmDelete ? (
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-marble-rust">Delete permanently?</span>
                    <button
                      type="button"
                      onClick={() => void handleDelete()}
                      className="rounded-lg bg-marble-rust px-3 py-1.5 text-xs font-semibold text-white"
                    >
                      Confirm
                    </button>
                    <button
                      type="button"
                      onClick={() => setConfirmDelete(false)}
                      className="text-xs text-slate-muted hover:text-ink"
                    >
                      Cancel
                    </button>
                  </div>
                ) : (
                  <button
                    type="button"
                    onClick={() => setConfirmDelete(true)}
                    className="rounded-xl border border-hairline px-3 py-1.5 text-xs font-medium text-slate-muted hover:border-marble-rust/50 hover:text-marble-rust"
                  >
                    Delete Grave
                  </button>
                )}
              </div>

              {idea.status !== 'resurrected' && (
                <button
                  type="button"
                  onClick={() => void handleResurrect()}
                  className="inline-flex items-center gap-1.5 rounded-xl bg-azure px-4 py-2 text-xs font-semibold text-white shadow-sm transition-all hover:bg-azure-deep"
                >
                  <BoltIcon className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {idea.status === 'buried' ? 'Move to Life Support' : 'Resurrect Idea'}
                  </span>
                </button>
              )}
            </div>
          )}
        </motion.div>
      )}
    </Modal>
  );
}
