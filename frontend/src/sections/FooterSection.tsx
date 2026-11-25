import { TombstoneIcon } from '../components/Icons';

export function FooterSection() {
  return (
    <footer className="border-t border-hairline bg-alabaster py-12">
      <div className="mx-auto flex max-w-6xl flex-col items-start justify-between gap-6 px-4 sm:flex-row sm:items-center sm:px-6">
        <div className="flex items-center gap-3">
          <span className="flex h-9 w-9 items-center justify-center rounded-xl border border-hairline bg-canvas text-azure shadow-sm">
            <TombstoneIcon className="h-4 w-4" aria-hidden="true" />
          </span>
          <div>
            <p className="font-display text-sm font-semibold text-ink">
              Startup Idea Graveyard
            </p>
            <p className="text-xs text-slate-muted">
              Built with Go, Gin, PostgreSQL, React, Tailwind &amp; Framer Motion.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-5 text-xs font-medium text-slate-muted">
          <a href="#top" className="transition-colors hover:text-ink">
            Back to top
          </a>
          <a href="#graveyard" className="transition-colors hover:text-ink">
            Graveyard
          </a>
          <a href="#autopsy" className="transition-colors hover:text-ink">
            Autopsy Report
          </a>
          <a
            href="/api/docs/index.html"
            target="_blank"
            rel="noreferrer"
            className="transition-colors hover:text-azure"
          >
            API Docs
          </a>
        </div>
      </div>
    </footer>
  );
}
