import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        crypt: '#0A0B0D',
        headstone: '#121418',
        'slate-mist': '#1B1F27',
        'stone-border': '#272C37',
        bone: '#EAE6DF',
        stone: '#8E929B',
        necro: '#7CEB4F',
        'necro-dim': '#4A9C2D',
        rust: '#D97757',
        amber: '#F59E0B',
      },
      fontFamily: {
        display: ['Fraunces', 'Georgia', 'serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      boxShadow: {
        tombstone: '0 20px 50px -12px rgba(0, 0, 0, 0.75)',
        'necro-glow': '0 0 28px -4px rgba(124, 235, 79, 0.28)',
      },
    },
  },
  plugins: [],
};

export default config;
