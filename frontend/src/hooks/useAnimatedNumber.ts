import { useEffect, useState } from 'react';
import { useReducedMotion } from 'framer-motion';

export function useAnimatedNumber(target: number, durationMs = 900): number {
  const shouldReduce = useReducedMotion();
  const [value, setValue] = useState(shouldReduce ? target : 0);

  useEffect(() => {
    if (shouldReduce) {
      setValue(target);
      return;
    }

    let frame = 0;
    const start = performance.now();
    const from = 0;

    const tick = (now: number) => {
      const progress = Math.min(1, (now - start) / durationMs);
      const eased = 1 - Math.pow(1 - progress, 3);
      setValue(Math.round(from + (target - from) * eased));
      if (progress < 1) {
        frame = requestAnimationFrame(tick);
      }
    };

    frame = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(frame);
  }, [target, durationMs, shouldReduce]);

  return value;
}
