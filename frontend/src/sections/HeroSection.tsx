import { motion, useReducedMotion, useScroll, useTransform } from 'framer-motion';
import { useRef } from 'react';
import { Link } from 'react-router-dom';
import { ShovelIcon } from '../components/Icons';
import { useAnimatedNumber } from '../hooks/useAnimatedNumber';
import { useAuth } from '../hooks/useAuth';
import { formatMoneyCents } from '../lib/format';
import type { InsightsSummary } from '../types/api';

interface HeroSectionProps {
  summary: InsightsSummary | null;
}

export function HeroSection({ summary }: HeroSectionProps) {
  const { isAuthenticated } = useAuth();
  const containerRef = useRef<HTMLElement>(null);
  const shouldReduce = useReducedMotion();

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start start', 'end start'],
  });
  const fogY = useTransform(scrollYProgress, [0, 1], [0, shouldReduce ? 0 : 90]);
  const moonY = useTransform(scrollYProgress, [0, 1], [0, shouldReduce ? 0 : 45]);

  const animatedCount = useAnimatedNumber(summary?.total_ideas ?? 0);
  const animatedDollars = useAnimatedNumber(Math.round((summary?.total_money_lost_cents ?? 0) / 100));

  return (
    <section
      id="top"
      ref={containerRef}
      className="relative overflow-hidden border-b border-stone-border/60 bg-gradient-to-b from-crypt via-headstone/70 to-crypt pb-20 pt-16 sm:pb-28 sm:pt-24"
    >
      {/* Flickering Moon */}
      <motion.div
        style={{ y: moonY }}
        aria-hidden="true"
        className="pointer-events-none absolute right-8 top-8 h-24 w-24 rounded-full bg-gradient-to-br from-bone/90 via-bone/40 to-transparent blur-[1px] animate-moon-flicker sm:right-20 sm:top-12 sm:h-32 sm:w-32"
      />

      {/* Drifting Fog Layers */}
      <motion.div
        style={{ y: fogY }}
        aria-hidden="true"
        className="pointer-events-none absolute inset-x-0 bottom-0 h-44 bg-gradient-to-t from-necro/5 via-slate-mist/25 to-transparent blur-2xl animate-fog-slow"
      />

      <div className="relative mx-auto max-w-6xl px-4 sm:px-6">
        <motion.div
          initial={shouldReduce ? false : { opacity: 0, y: 18 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.55, ease: 'easeOut' }}
          className="max-w-3xl"
        >
          <p className="mb-4 inline-flex items-center gap-2 rounded-full border border-stone-border bg-headstone/90 px-3 py-1 font-mono text-xs text-necro">
            <span className="h-2 w-2 rounded-full bg-necro" />
            <span>
              {isAuthenticated ? 'Your Personal Crypt' : 'Viewing Demo Graveyard — Sign up to bury yours'}
            </span>
          </p>

          <h1 className="font-display text-4xl font-bold tracking-tight text-bone sm:text-6xl sm:leading-[1.08]">
            Startup Idea Graveyard
          </h1>

          <p className="mt-5 max-w-2xl text-base leading-relaxed text-stone sm:text-lg">
            Bury your abandoned side projects with an honest cause of death and an epitaph.
            Then inspect the autopsy report to see why your ideas keep dying.
          </p>

          <div className="mt-8 flex flex-wrap items-center gap-4">
            <Link
              to={isAuthenticated ? '/bury' : '/register'}
              className="inline-flex items-center gap-2 rounded-xl bg-necro px-6 py-3 text-sm font-semibold text-crypt shadow-necro-glow transition-transform hover:scale-[1.02] active:scale-[0.98]"
            >
              <ShovelIcon className="h-4 w-4" aria-hidden="true" />
              <span>{isAuthenticated ? 'Bury an Idea' : 'Start Your Graveyard'}</span>
            </Link>

            <a
              href="#autopsy"
              className="inline-flex items-center rounded-xl border border-stone-border bg-headstone px-5 py-3 text-sm font-medium text-bone transition-colors hover:border-necro/60"
            >
              Read the Autopsy
            </a>
          </div>
        </motion.div>

        {/* Live Counters */}
        <motion.div
          initial={shouldReduce ? false : { opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.55, delay: 0.15, ease: 'easeOut' }}
          className="mt-12 grid max-w-2xl grid-cols-2 gap-4 sm:grid-cols-3"
        >
          <div className="rounded-xl border border-stone-border bg-headstone/80 p-4">
            <div className="font-mono text-2xl font-semibold text-necro sm:text-3xl">
              {animatedCount}
            </div>
            <div className="mt-1 text-xs uppercase tracking-wider text-stone">
              Ideas Buried
            </div>
          </div>

          <div className="rounded-xl border border-stone-border bg-headstone/80 p-4">
            <div className="font-mono text-2xl font-semibold text-bone sm:text-3xl">
              {formatMoneyCents(animatedDollars * 100)}
            </div>
            <div className="mt-1 text-xs uppercase tracking-wider text-stone">
              Capital Burned
            </div>
          </div>

          <div className="col-span-2 rounded-xl border border-stone-border bg-headstone/80 p-4 sm:col-span-1">
            <div className="font-mono text-2xl font-semibold text-bone sm:text-3xl">
              {summary ? `${Math.round(summary.avg_lifespan_days)}d` : '—'}
            </div>
            <div className="mt-1 text-xs uppercase tracking-wider text-stone">
              Avg Lifespan
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
