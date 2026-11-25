import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ApiError } from '../api/client';
import { Modal } from '../components/Modal';
import { useAuth } from '../hooks/useAuth';

interface AuthModalProps {
  mode: 'login' | 'register';
}

export function AuthModal({ mode }: AuthModalProps) {
  const { login, register } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const close = () => navigate('/');

  const validate = (): boolean => {
    const errs: Record<string, string> = {};
    if (!email.trim() || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) {
      errs.email = 'Enter a valid email address.';
    }
    if (mode === 'register' && password.length < 8) {
      errs.password = 'Password must be at least 8 characters.';
    } else if (mode === 'login' && !password) {
      errs.password = 'Password is required.';
    }
    if (mode === 'register' && displayName.trim().length < 2) {
      errs.displayName = 'Display name must be at least 2 characters.';
    }
    setFieldErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setServerError(null);
    if (!validate()) return;

    setSubmitting(true);
    try {
      if (mode === 'login') {
        await login(email.trim(), password);
      } else {
        await register(email.trim(), password, displayName.trim());
      }
      navigate('/');
    } catch (err) {
      if (err instanceof ApiError) {
        setServerError(err.message);
      } else {
        setServerError('Something went wrong. Please try again.');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal
      title={mode === 'login' ? 'Enter Your Crypt' : 'Claim Your Graveyard Plot'}
      titleId="auth-modal-title"
      onClose={close}
      maxWidthClass="max-w-md"
    >
      <form onSubmit={(e) => void onSubmit(e)} noValidate className="space-y-4">
        {serverError && (
          <div role="alert" className="rounded-lg border border-rust/50 bg-rust/10 px-3 py-2 text-xs text-rust">
            {serverError}
          </div>
        )}

        {mode === 'register' && (
          <div>
            <label htmlFor="auth-name" className="block text-xs font-medium text-bone">
              Display Name
            </label>
            <input
              id="auth-name"
              type="text"
              autoComplete="name"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              placeholder="Serial Pivoter"
              className="mt-1.5 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone placeholder:text-stone/50 focus:border-necro focus:outline-none"
            />
            {fieldErrors.displayName && (
              <p className="mt-1 text-xs text-rust">{fieldErrors.displayName}</p>
            )}
          </div>
        )}

        <div>
          <label htmlFor="auth-email" className="block text-xs font-medium text-bone">
            Email
          </label>
          <input
            id="auth-email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="founder@example.com"
            className="mt-1.5 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone placeholder:text-stone/50 focus:border-necro focus:outline-none"
          />
          {fieldErrors.email && (
            <p className="mt-1 text-xs text-rust">{fieldErrors.email}</p>
          )}
        </div>

        <div>
          <label htmlFor="auth-password" className="block text-xs font-medium text-bone">
            Password
          </label>
          <input
            id="auth-password"
            type="password"
            autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            className="mt-1.5 w-full rounded-lg border border-stone-border bg-crypt px-3 py-2 text-sm text-bone placeholder:text-stone/50 focus:border-necro focus:outline-none"
          />
          {fieldErrors.password && (
            <p className="mt-1 text-xs text-rust">{fieldErrors.password}</p>
          )}
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full rounded-xl bg-necro py-2.5 text-sm font-semibold text-crypt transition-opacity hover:opacity-95 disabled:opacity-50"
        >
          {submitting
            ? 'Opening gates...'
            : mode === 'login'
              ? 'Sign In'
              : 'Create Account'}
        </button>

        <p className="text-center text-xs text-stone">
          {mode === 'login' ? (
            <>
              New here?{' '}
              <Link to="/register" className="text-necro hover:underline">
                Claim a plot
              </Link>
            </>
          ) : (
            <>
              Already have a crypt?{' '}
              <Link to="/login" className="text-necro hover:underline">
                Sign in
              </Link>
            </>
          )}
        </p>
      </form>
    </Modal>
  );
}
