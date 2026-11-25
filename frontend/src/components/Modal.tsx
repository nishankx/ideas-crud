import { motion, useReducedMotion } from 'framer-motion';
import { useEffect, useRef, type ReactNode } from 'react';
import { CloseIcon } from './Icons';

interface ModalProps {
  title: string;
  titleId: string;
  onClose: () => void;
  children: ReactNode;
  maxWidthClass?: string;
}

export function Modal({
  title,
  titleId,
  onClose,
  children,
  maxWidthClass = 'max-w-lg',
}: ModalProps) {
  const panelRef = useRef<HTMLDivElement>(null);
  const shouldReduce = useReducedMotion();

  useEffect(() => {
    const previousFocus = document.activeElement as HTMLElement | null;
    const panel = panelRef.current;

    const focusables = panel?.querySelectorAll<HTMLElement>(
      'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])',
    );
    if (focusables && focusables.length > 0) {
      focusables[0].focus();
    } else {
      panel?.focus();
    }

    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault();
        onClose();
        return;
      }
      if (e.key === 'Tab' && panel) {
        const items = panel.querySelectorAll<HTMLElement>(
          'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])',
        );
        if (items.length === 0) return;
        const first = items[0];
        const last = items[items.length - 1];
        if (e.shiftKey && document.activeElement === first) {
          e.preventDefault();
          last.focus();
        } else if (!e.shiftKey && document.activeElement === last) {
          e.preventDefault();
          first.focus();
        }
      }
    };

    document.addEventListener('keydown', onKeyDown);
    return () => {
      document.removeEventListener('keydown', onKeyDown);
      previousFocus?.focus?.();
    };
  }, [onClose]);

  return (
    <motion.div
      initial={shouldReduce ? false : { opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={shouldReduce ? undefined : { opacity: 0 }}
      transition={{ duration: 0.2 }}
      className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-black/75 p-4 backdrop-blur-sm"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <motion.div
        ref={panelRef}
        tabIndex={-1}
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        initial={shouldReduce ? false : { opacity: 0, scale: 0.96, y: 12 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        exit={shouldReduce ? undefined : { opacity: 0, scale: 0.96, y: 12 }}
        transition={{ type: 'spring', damping: 26, stiffness: 320 }}
        className={`relative my-8 w-full ${maxWidthClass} rounded-2xl border border-stone-border bg-headstone p-6 shadow-tombstone`}
      >
        <div className="flex items-start justify-between gap-4 border-b border-stone-border/70 pb-4">
          <h2 id={titleId} className="font-display text-xl font-bold text-bone">
            {title}
          </h2>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close modal"
            className="rounded-lg border border-stone-border bg-crypt p-1.5 text-stone transition-colors hover:border-necro/60 hover:text-bone"
          >
            <CloseIcon className="h-4 w-4" aria-hidden="true" />
          </button>
        </div>

        <div className="mt-5">{children}</div>
      </motion.div>
    </motion.div>
  );
}
