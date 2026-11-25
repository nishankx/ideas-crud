import { motion, useReducedMotion, useScroll, useSpring, useTransform } from 'framer-motion';
import { useRef } from 'react';
import { Link } from 'react-router-dom';
import { ShovelIcon, TombstoneIcon } from '../components/Icons';
import { useAnimatedNumber } from '../hooks/useAnimatedNumber';
import { useAuth } from '../hooks/useAuth';
import { CAUSE_LABELS, formatMoneyCents } from '../lib/format';
import type { InsightsSummary } from '../types/api';

interface HeroSectionProps {
  summary: InsightsSummary | null;
}

const SMOOTH_EASE = [0.22, 1, 0.36, 1] as const;

export function HeroSection({ summary }: HeroSectionProps) {
  const { isAuthenticated } = useAuth();
  const containerRef = useRef<HTMLElement>(null);
  const shouldReduce = useReducedMotion();

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ['start start', 'end start'],
  });
  const smoothProgress = useSpring(scrollYProgress, {
    stiffness: 90,
    damping: 26,
    restDelta: 0.001,
  });

  const cloudBackY = useTransform(smoothProgress, [0, 1], [0, shouldReduce ? 0 : 48]);
  const cloudFrontY = useTransform(smoothProgress, [0, 1], [0, shouldReduce ? 0 : -28]);
  const portalScale = useTransform(smoothProgress, [0, 1], [1, shouldReduce ? 1 : 0.94]);

  const animatedCount = useAnimatedNumber(summary?.total_ideas ?? 0);
  const animatedDollars = useAnimatedNumber(Math.round((summary?.total_money_lost_cents ?? 0) / 100));

  const marqueeCauses = Object.values(CAUSE_LABELS);

  return (
    <section id="top" ref={containerRef} className="bg-canvas pb-12 pt-5 sm:pb-16 sm:pt-7">
      {/* Sky & Clouds Surrealist Banner */}
      <div className="mx-auto max-w-6xl px-4 sm:px-6">
        <motion.div
          initial={shouldReduce ? false : { opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, ease: SMOOTH_EASE }}
          className="relative h-64 overflow-hidden rounded-3xl bg-gradient-to-b from-[#1B75E8] via-[#4CA2F7] to-[#98CBFD] shadow-tombstone sm:h-80"
        >
          {/* Distant Sunlit Cloud Layer */}
          <motion.svg
            style={{ y: cloudBackY }}
            viewBox="0 0 1200 400"
            aria-hidden="true"
            className="pointer-events-none absolute inset-0 h-full w-full object-cover opacity-85 animate-cloud-float-slow"
            preserveAspectRatio="xMidYMid slice"
          >
            <g fill="rgba(255,255,255,0.68)">
              <circle cx="140" cy="240" r="110" />
              <circle cx="240" cy="200" r="130" />
              <circle cx="360" cy="250" r="105" />
              <circle cx="860" cy="220" r="120" />
              <circle cx="990" cy="185" r="145" />
              <circle cx="1110" cy="245" r="115" />
            </g>
          </motion.svg>

          {/* Surrealist Torn-Cloud Portal Centerpiece */}
          <motion.div
            style={{ scale: portalScale }}
            initial={shouldReduce ? false : { opacity: 0, scale: 0.92 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, delay: 0.12, ease: SMOOTH_EASE }}
            className="relative z-10 flex h-full items-center justify-center"
          >
            <div className="relative flex h-40 w-44 flex-col items-center justify-center rounded-t-[4rem] rounded-b-2xl border-4 border-white/95 bg-gradient-to-b from-[#1F2430] to-[#0D1017] px-4 text-center shadow-2xl sm:h-48 sm:w-52">
              <div className="absolute -inset-2 -z-10 rounded-t-[4.5rem] rounded-b-3xl bg-white/40 blur-md" />
              <TombstoneIcon className="h-8 w-8 text-azure sm:h-9 sm:w-9" aria-hidden="true" />
              <span className="mt-2 font-mono text-[10px] uppercase tracking-[0.22em] text-white/60">
                R.I.P.
              </span>
              <span className="mt-1 font-display text-sm font-medium italic text-white/95 sm:text-base">
                “It scaled in our heads.”
              </span>
            </div>
          </motion.div>

          {/* Foreground Fluffy Cloud Bank */}
          <motion.svg
            style={{ y: cloudFrontY }}
            viewBox="0 0 1200 200"
            aria-hidden="true"
            className="pointer-events-none absolute -bottom-6 inset-x-0 z-20 w-full animate-cloud-float"
            preserveAspectRatio="none"
          >
            <g fill="#FFFFFF">
              <circle cx="70" cy="170" r="85" />
              <circle cx="190" cy="185" r="95" />
              <circle cx="330" cy="195" r="80" />
              <circle cx="480" cy="205" r="75" />
              <circle cx="730" cy="205" r="75" />
              <circle cx="880" cy="190" r="88" />
              <circle cx="1020" cy="175" r="98" />
              <circle cx="1150" cy="175" r="85" />
            </g>
          </motion.svg>
        </motion.div>
      </div>

      {/* Centered Editorial Headline */}
      <div className="mx-auto max-w-4xl px-4 pt-10 text-center sm:px-6 sm:pt-12">
        <motion.div
          initial={shouldReduce ? false : { opacity: 0, y: 14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.65, delay: 0.15, ease: SMOOTH_EASE }}
          className="inline-flex items-center gap-2 rounded-full border border-hairline bg-alabaster px-3.5 py-1 text-xs font-medium text-slate-muted"
        >
          <span className="h-2 w-2 rounded-full bg-azure" />
          <span>
            {isAuthenticated
              ? 'Personal Crypt Active'
              : 'Viewing Demo Graveyard — Sign up to bury yours'}
          </span>
        </motion.div>

        <motion.h1
          initial={shouldReduce ? false : { opacity: 0, y: 18 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.22, ease: SMOOTH_EASE }}
          className="mt-5 font-display text-4xl font-medium tracking-[-0.025em] text-ink sm:text-6xl sm:leading-[1.06]"
        >
          Startup Idea Graveyard
          <br />
          <span className="text-slate-muted">for Serial Builders</span>
        </motion.h1>

        <motion.p
          initial={shouldReduce ? false : { opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.65, delay: 0.3, ease: SMOOTH_EASE }}
          className="mx-auto mt-4 max-w-xl text-sm leading-relaxed text-slate-muted sm:text-base"
        >
          Bury abandoned projects with an honest cause of death and an epitaph, then inspect the autopsy report to see why your ideas keep dying.
        </motion.p>

        <motion.div
          initial={shouldReduce ? false : { opacity: 0, y: 14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.65, delay: 0.38, ease: SMOOTH_EASE }}
          className="mt-7 flex flex-wrap items-center justify-center gap-3"
        >
          <Link
            to={isAuthenticated ? '/bury' : '/register'}
            className="inline-flex items-center gap-2 rounded-xl bg-azure px-6 py-3 text-sm font-semibold text-white shadow-azure-glow transition-all hover:bg-azure-deep active:scale-[0.98]"
          >
            <ShovelIcon className="h-4 w-4" aria-hidden="true" />
            <span>{isAuthenticated ? 'Bury an Idea' : 'Start Your Graveyard'}</span>
          </Link>

          <a
            href="#autopsy"
            className="inline-flex items-center rounded-xl border border-hairline bg-canvas px-5 py-3 text-sm font-medium text-ink transition-colors hover:border-ink/30 hover:bg-alabaster"
          >
            Read the Autopsy
          </a>
        </motion.div>

        {/* Live Counters */}
        <motion.div
          initial={shouldReduce ? false : { opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.65, delay: 0.45, ease: SMOOTH_EASE }}
          className="mx-auto mt-10 grid max-w-xl grid-cols-3 divide-x divide-hairline rounded-2xl border border-hairline bg-alabaster/70 py-4"
        >
          <div className="px-3">
            <div className="font-mono text-xl font-semibold text-azure sm:text-2xl">
              {animatedCount}
            </div>
            <div className="mt-0.5 text-[11px] font-medium uppercase tracking-wider text-slate-muted">
              Buried
            </div>
          </div>
          <div className="px-3">
            <div className="font-mono text-xl font-semibold text-ink sm:text-2xl">
              {formatMoneyCents(animatedDollars * 100)}
            </div>
            <div className="mt-0.5 text-[11px] font-medium uppercase tracking-wider text-slate-muted">
              Burned
            </div>
          </div>
          <div className="px-3">
            <div className="font-mono text-xl font-semibold text-ink sm:text-2xl">
              {summary ? `${Math.round(summary.avg_lifespan_days)}d` : '—'}
            </div>
            <div className="mt-0.5 text-[11px] font-medium uppercase tracking-wider text-slate-muted">
              Avg Life
            </div>
          </div>
        </motion.div>
      </div>

      {/* Horizontal Marquee Strip (matching reference ticker) */}
      <div className="mt-12 overflow-hidden border-y border-hairline bg-alabaster/60 py-3.5">
        <div className="flex w-max items-center gap-10 whitespace-nowrap animate-marquee">
          {[...marqueeCauses, ...marqueeCauses].map((label, idx) => (
            <span
              key={`${label}-${idx}`}
              className="inline-flex items-center gap-2 font-mono text-xs font-medium uppercase tracking-widest text-slate-muted/80"
            >
              <span className="h-1.5 w-1.5 rounded-full bg-azure/50" />
              {label}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}
