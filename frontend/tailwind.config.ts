import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        canvas: '#FFFFFF',
        alabaster: '#F7F9FC',
        'cloud-wash': '#EEF5FE',
        hairline: '#E4E7EC',
        ink: '#111827',
        'slate-muted': '#64748B',
        azure: '#2589F5',
        'azure-deep': '#196FD3',
        'marble-rust': '#C85A32',
        amber: '#D97706',
        revive: '#059669',
      },
      fontFamily: {
        display: ['Fraunces', 'Georgia', 'serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      boxShadow: {
        tombstone: '0 14px 36px -14px rgba(17, 24, 39, 0.10), 0 2px 8px -2px rgba(17, 24, 39, 0.05)',
        'azure-glow': '0 16px 36px -10px rgba(37, 137, 245, 0.24)',
        pill: '0 4px 20px -4px rgba(17, 24, 39, 0.08)',
      },
    },
  },
  plugins: [],
};

export default config;
