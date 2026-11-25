import { motion, useReducedMotion } from 'framer-motion';
import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { ApiError, apiRequest } from '../api/client';
import { ShovelIcon } from '../components/Icons';
import { Modal } from '../components/Modal';
import { useAuth } from '../hooks/useAuth';
import { CATEGORY_LABELS, CAUSE_LABELS, VALIDATION_LABELS } from '../lib/format';
import type { DeathCause, Idea, IdeaCategory, ValidationLevel } from '../types/api';

interface BuryModalProps {
  onBuried: () => void;
}

export function BuryModal({ onBuried }: BuryModalProps) {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const shouldReduce = useReducedMotion();

  const today = new Date().toISOString().slice(0, 10);

  const [name, setName] = useState('');
  const [tagline, setTagline] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState<IdeaCategory>('saas');
  const [bornAt, setBornAt] = useState(today);
  const [diedAt, setDiedAt] = useState(today);
  const [primaryCause, setPrimaryCause] = useState<DeathCause>('no_market_need');
  const [secondaryCauses, setSecondaryCauses] = useState<DeathCause[]>([]);
  const [epitaph, setEpitaph] = useState('');
  const [autopsy, setAutopsy] = useState('');
  const [validationLevel, setValidationLevel] = useState<ValidationLevel>('none');
  const [moneyDollars, setMoneyDollars] = useState('0');
  const [hoursInvested, setHoursInvested] = useState('0');
  const [wouldResurrect, setWouldResurrect] = useState(false);

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [buryingAnim, setBuryingAnim] = useState(false);

  const close = () => navigate('/');

  if (!isAuthenticated) {
    return (
      <Modal title="Sign In to Bury an Idea" titleId="bury-auth-title" onClose={close}>
        <p className="text-sm text-stone">
          You are currently browsing the public demo graveyard. Create a free account to bury and analyze your own abandoned ideas.
        </p>
        <div className="mt-5 flex justify-end gap-3">
          <button
            type="button"
            onClick={close}
            className="rounded-lg border border-stone-border px-4 py-2 text-xs text-stone hover:text-bone"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={() => navigate('/register')}
            className="rounded-lg bg-necro px-4 py-2 text-xs font-semibold text-crypt"
          >
            Create Account
          </button>
        </div>
      </Modal>
    );
  }

  const toggleSecondary = (cause: DeathCause) => {
    if (cause === primaryCause) return;
    setSecondaryCauses((prev) =>
      prev.includes(cause) ? prev.filter((c) => c !== cause) : [...prev, cause],
    );
  };

  const validate = (): boolean => {
    const errs: Record<string, string> = {};
    if (!name.trim()) errs.name = 'Name your fallen startup.';
    if (!epitaph.trim()) errs.epitaph = 'Every tombstone needs an epitaph.';
    if (epitaph.trim().length > 140) errs.epitaph = 'Epitaph must be 140 characters or fewer.';
    if (diedAt < bornAt) errs.diedAt = 'Date of death cannot be before birth.';
    const dollarsNum = Number(moneyDollars);
    if (Number.isNaN(dollarsNum) || dollarsNum < 0) errs.moneyDollars = 'Enter a non-negative dollar amount.';
    const hoursNum = Number(hoursInvested);
    if (Number.isNaN(hoursNum) || hoursNum < 0 || !Number.isInteger(hoursNum)) {
      errs.hoursInvested = 'Enter a whole number of hours.';
    }
    setErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setServerError(null);
    if (!validate()) return;

    setBuryingAnim(true);
    try {
      await apiRequest<Idea>('/ideas', {
        method: 'POST',
        body: JSON.stringify({
          name: name.trim(),
          tagline: tagline.trim(),
          description: description.trim(),
          category,
          born_at: bornAt,
          died_at: diedAt,
          primary_cause: primaryCause,
          secondary_causes: secondaryCauses,
          epitaph: epitaph.trim(),
          autopsy: autopsy.trim(),
          validation_level: validationLevel,
          money_spent_cents: Math.round(Number(moneyDollars) * 100),
          hours_invested: Number(hoursInvested),
          would_resurrect: wouldResurrect,
        }),
      });

      if (!shouldReduce) {
        await new Promise((r) => setTimeout(r, 650));
      }
      onBuried();
      navigate('/');
    } catch (err) {
      setBuryingAnim(false);
      if (err instanceof ApiError) {
        setServerError(err.message);
      } else {
        setServerError('Failed to bury idea. Please try again.');
      }
    }
  };

  return (
    <Modal title="Bury an Idea" titleId="bury-modal-title" onClose={close} maxWidthClass="max-w-2xl">
      {buryingAnim ? (
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <motion.div
            initial={shouldReduce ? false : { rotate: -25, y: -12 }}
            animate={{ rotate: [0, -35, 10, 0], y: [0, 8, 0] }}
            transition={{ duration: 0.6 }}
            className="rounded-full border border-necro/40 bg-necro/10 p-5 text-necro"
          >
            <ShovelIcon className="h-10 w-10" aria-hidden="true" />
          </motion.div>
          <p className="mt-4 font-display text-lg font-semibold text-bone">
            Lowering {name} into the dirt...
          </p>
          <p className="mt-1 text-xs text-stone">May its lessons live on.</p>
        </div>
      ) : (
        <form onSubmit={(e) => void onSubmit(e)} noValidate className="space-y-4">
          {serverError && (
            <div role="alert" className="rounded-lg border border-rust/50 bg-rust/10 px-3 py-2 text-xs text-rust">
              {serverError}
            </div>
          )}

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label htmlFor="bury-name" className="block text-xs font-medium text-bone">
                Startup Name *
              </label>
              <input
                id="bury-name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="BarkCRM"
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
              />
              {errors.name && <p className="mt-1 text-xs text-rust">{errors.name}</p>}
            </div>

            <div>
              <label htmlFor="bury-category" className="block text-xs font-medium text-bone">
                Category *
              </label>
              <select
                id="bury-category"
                value={category}
                onChange={(e) => setCategory(e.target.value as IdeaCategory)}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
              >
                {Object.entries(CATEGORY_LABELS).map(([k, label]) => (
                  <option key={k} value={k}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div>
            <label htmlFor="bury-tagline" className="block text-xs font-medium text-bone">
              One-line Pitch
            </label>
            <input
              id="bury-tagline"
              type="text"
              value={tagline}
              onChange={(e) => setTagline(e.target.value)}
              placeholder="Salesforce for neighborhood dog walkers"
              className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
            />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label htmlFor="bury-born" className="block text-xs font-medium text-bone">
                Born Date *
              </label>
              <input
                id="bury-born"
                type="date"
                value={bornAt}
                onChange={(e) => setBornAt(e.target.value)}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 font-mono text-sm text-bone focus:border-necro focus:outline-none"
              />
            </div>

            <div>
              <label htmlFor="bury-died" className="block text-xs font-medium text-bone">
                Died Date *
              </label>
              <input
                id="bury-died"
                type="date"
                value={diedAt}
                onChange={(e) => setDiedAt(e.target.value)}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 font-mono text-sm text-bone focus:border-necro focus:outline-none"
              />
              {errors.diedAt && <p className="mt-1 text-xs text-rust">{errors.diedAt}</p>}
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label htmlFor="bury-cause" className="block text-xs font-medium text-bone">
                Primary Cause of Death *
              </label>
              <select
                id="bury-cause"
                value={primaryCause}
                onChange={(e) => {
                  const next = e.target.value as DeathCause;
                  setPrimaryCause(next);
                  setSecondaryCauses((prev) => prev.filter((c) => c !== next));
                }}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
              >
                {Object.entries(CAUSE_LABELS).map(([k, label]) => (
                  <option key={k} value={k}>
                    {label}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="bury-validation" className="block text-xs font-medium text-bone">
                Validation Reached
              </label>
              <select
                id="bury-validation"
                value={validationLevel}
                onChange={(e) => setValidationLevel(e.target.value as ValidationLevel)}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
              >
                {Object.entries(VALIDATION_LABELS).map(([k, label]) => (
                  <option key={k} value={k}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div>
            <span className="block text-xs font-medium text-bone">
              Secondary Causes (optional)
            </span>
            <div className="mt-1.5 flex flex-wrap gap-1.5">
              {Object.entries(CAUSE_LABELS)
                .filter(([k]) => k !== primaryCause)
                .map(([k, label]) => {
                  const active = secondaryCauses.includes(k as DeathCause);
                  return (
                    <button
                      key={k}
                      type="button"
                      onClick={() => toggleSecondary(k as DeathCause)}
                      className={`rounded-full border px-2.5 py-1 text-[11px] transition-colors ${
                        active
                          ? 'border-necro bg-necro/15 text-necro'
                          : 'border-stone-border bg-crypt text-stone hover:text-bone'
                      }`}
                    >
                      {label}
                    </button>
                  );
                })}
            </div>
          </div>

          <div>
            <div className="flex items-center justify-between">
              <label htmlFor="bury-epitaph" className="block text-xs font-medium text-bone">
                Epitaph (engraved on the stone) *
              </label>
              <span className={`font-mono text-[11px] ${epitaph.length > 140 ? 'text-rust' : 'text-stone'}`}>
                {epitaph.length}/140
              </span>
            </div>
            <input
              id="bury-epitaph"
              type="text"
              maxLength={140}
              value={epitaph}
              onChange={(e) => setEpitaph(e.target.value)}
              placeholder="Here lies a CRM for dog walkers. It never met a customer."
              className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 font-display text-sm italic text-bone focus:border-necro focus:outline-none"
            />
            {errors.epitaph && <p className="mt-1 text-xs text-rust">{errors.epitaph}</p>}
          </div>

          <div>
            <label htmlFor="bury-autopsy" className="block text-xs font-medium text-bone">
              Autopsy (what did you learn?)
            </label>
            <textarea
              id="bury-autopsy"
              rows={3}
              value={autopsy}
              onChange={(e) => setAutopsy(e.target.value)}
              placeholder="Talk to 10 customers before buying the domain name."
              className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
            />
          </div>

          <div>
            <label htmlFor="bury-description" className="block text-xs font-medium text-bone">
              Description (optional)
            </label>
            <textarea
              id="bury-description"
              rows={2}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone focus:border-necro focus:outline-none"
            />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label htmlFor="bury-money" className="block text-xs font-medium text-bone">
                Money Lost (USD)
              </label>
              <input
                id="bury-money"
                type="number"
                min={0}
                step="1"
                value={moneyDollars}
                onChange={(e) => setMoneyDollars(e.target.value)}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 font-mono text-sm text-bone focus:border-necro focus:outline-none"
              />
              {errors.moneyDollars && <p className="mt-1 text-xs text-rust">{errors.moneyDollars}</p>}
            </div>

            <div>
              <label htmlFor="bury-hours" className="block text-xs font-medium text-bone">
                Hours Invested
              </label>
              <input
                id="bury-hours"
                type="number"
                min={0}
                step="1"
                value={hoursInvested}
                onChange={(e) => setHoursInvested(e.target.value)}
                className="mt-1 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 font-mono text-sm text-bone focus:border-necro focus:outline-none"
              />
              {errors.hoursInvested && <p className="mt-1 text-xs text-rust">{errors.hoursInvested}</p>}
            </div>
          </div>

          <label className="flex items-center gap-2.5 text-xs text-stone">
            <input
              type="checkbox"
              checked={wouldResurrect}
              onChange={(e) => setWouldResurrect(e.target.checked)}
              className="h-4 w-4 rounded border-stone-border bg-crypt text-necro focus:ring-necro"
            />
            <span>I would consider resurrecting this idea someday</span>
          </label>

          <div className="flex justify-end gap-3 border-t border-stone-border/70 pt-4">
            <button
              type="button"
              onClick={close}
              className="rounded-lg border border-stone-border px-4 py-2 text-xs text-stone hover:text-bone"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="inline-flex items-center gap-1.5 rounded-lg bg-necro px-5 py-2 text-xs font-semibold text-crypt"
            >
              <ShovelIcon className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Bury Idea</span>
            </button>
          </div>
        </form>
      )}
    </Modal>
  );
}
