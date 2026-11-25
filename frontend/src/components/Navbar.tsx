import { Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { ShovelIcon, TombstoneIcon } from './Icons';

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();

  return (
    <header className="sticky top-0 z-40 border-b border-hairline/80 bg-canvas/90 backdrop-blur-md">
      <nav aria-label="Primary" className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3 sm:px-6">
        <div className="flex items-center gap-4">
          <a
            href="#top"
            className="group flex items-center gap-2.5 rounded-xl border border-hairline bg-canvas px-3 py-1.5 text-ink shadow-pill transition-colors hover:border-azure/40"
          >
            <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-cloud-wash text-azure">
              <TombstoneIcon className="h-4 w-4" aria-hidden="true" />
            </span>
            <span className="font-display text-base font-semibold tracking-tight">
              Graveyard
            </span>
          </a>

          <div className="hidden items-center gap-5 text-sm font-medium text-slate-muted md:flex">
            <a href="#graveyard" className="transition-colors hover:text-ink">
              The Graves
            </a>
            <a href="#autopsy" className="transition-colors hover:text-ink">
              Autopsy Report
            </a>
            {!isAuthenticated && (
              <span className="rounded-full border border-azure/25 bg-cloud-wash px-2.5 py-0.5 font-mono text-[11px] font-medium text-azure-deep">
                Demo Mode
              </span>
            )}
          </div>
        </div>

        <div className="flex items-center gap-2.5">
          {isAuthenticated ? (
            <>
              <span className="hidden text-xs font-medium text-slate-muted sm:inline">
                {user?.display_name}
              </span>
              <button
                type="button"
                onClick={() => void logout()}
                className="rounded-lg border border-hairline bg-canvas px-3 py-1.5 text-xs font-medium text-slate-muted transition-colors hover:border-ink/30 hover:text-ink"
              >
                Log out
              </button>
            </>
          ) : (
            <Link
              to="/login"
              className="rounded-lg px-3 py-1.5 text-xs font-medium text-slate-muted transition-colors hover:text-ink"
            >
              Sign in
            </Link>
          )}

          <Link
            to={isAuthenticated ? '/bury' : '/register'}
            className="inline-flex items-center gap-1.5 rounded-lg bg-azure px-3.5 py-2 text-xs font-semibold text-white shadow-sm transition-all hover:bg-azure-deep active:scale-[0.98]"
          >
            <ShovelIcon className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Bury an Idea</span>
          </Link>
        </div>
      </nav>
    </header>
  );
}
