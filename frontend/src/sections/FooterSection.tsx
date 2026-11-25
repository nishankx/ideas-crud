import { TombstoneIcon } from '../components/Icons';

export function FooterSection() {
  return (
    <footer className="border-t border-stone-border/70 bg-headstone/60 py-12">
      <div className="mx-auto flex max-w-6xl flex-col items-start justify-between gap-6 px-4 sm:flex-row sm:items-center sm:px-6">
        <div className="flex items-center gap-3">
          <TombstoneIcon className="h-5 w-5 text-necro" aria-hidden="true" />
          <div>
            <p className="font-display text-sm font-semibold text-bone">
              Startup Idea Graveyard
            </p>
            <p className="text-xs text-stone">
              Built with Go, Gin, PostgreSQL, React, Tailwind &amp; Framer Motion.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-5 text-xs text-stone">
          <a href="#top" className="hover:text-bone">
            Back to top
          </a>
          <a href="#graveyard" className="hover:text-bone">
            Graveyard
          </a>
          <a href="#autopsy" className="hover:text-bone">
            Autopsy Report
          </a>
          <a
            href="/api/docs/index.html"
            target="_blank"
            rel="noreferrer"
            className="hover:text-necro"
          >
            API Docs
          </a>
        </div>
      </div>
    </footer>
  );
}
