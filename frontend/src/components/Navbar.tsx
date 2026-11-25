import { Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { ShovelIcon, TombstoneIcon } from './Icons';

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();

  return (
    <header className="sticky top-0 z-40 border-b border-stone-border/70 bg-crypt/85 backdrop-blur-md">
      <nav aria-label="Primary" className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3 sm:px-6">
        <a href="#top" className="group flex items-center gap-2.5 text-bone">
          <span className="flex h-9 w-9 items-center justify-center rounded-lg border border-stone-border bg-headstone text-necro transition-colors group-hover:border-necro/60">
            <TombstoneIcon className="h-5 w-5" aria-hidden="true" />
          </span>
          <span className="font-display text-lg font-semibold tracking-tight">
            Idea Graveyard
          </span>
        </a>

        <div className="hidden items-center gap-6 text-sm text-stone md:flex">
          <a href="#graveyard" className="transition-colors hover:text-bone">
            The Graveyard
          </a>
          <a href="#autopsy" className="transition-colors hover:text-bone">
            Autopsy Report
          </a>
          {!isAuthenticated && (
            <span className="rounded-full border border-stone-border bg-headstone px-2.5 py-0.5 font-mono text-xs text-necro">
              Demo Mode
            </span>
          )}
        </div>

        <div className="flex items-center gap-2.5">
          {isAuthenticated ? (
            <>
              <span className="hidden text-xs text-stone sm:inline">
                {user?.display_name}
              </span>
              <button
                type="button"
                onClick={() => void logout()}
                className="rounded-lg border border-stone-border bg-headstone px-3 py-1.5 text-xs font-medium text-stone transition-colors hover:border-stone hover:text-bone"
              >
                Log out
              </button>
            </>
          ) : (
            <Link
              to="/login"
              className="rounded-lg border border-stone-border bg-headstone px-3 py-1.5 text-xs font-medium text-bone transition-colors hover:border-necro/60"
            >
              Sign in
            </Link>
          )}

          <Link
            to={isAuthenticated ? '/bury' : '/register'}
            className="inline-flex items-center gap-1.5 rounded-lg bg-necro px-3.5 py-1.5 text-xs font-semibold text-crypt transition-transform hover:scale-[1.02] active:scale-[0.98]"
          >
            <ShovelIcon className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Bury an Idea</span>
          </Link>
        </div>
      </nav>
    </header>
  );
}
