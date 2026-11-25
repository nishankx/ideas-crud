import { MotionConfig } from 'framer-motion';
import { useState } from 'react';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { Navbar } from './components/Navbar';
import { AuthProvider, useAuth } from './hooks/useAuth';
import { useFetch } from './hooks/useFetch';
import { FooterSection } from './sections/FooterSection';
import { GraveyardSection } from './sections/GraveyardSection';
import { HeroSection } from './sections/HeroSection';
import type { InsightsSummary } from './types/api';

function GraveyardPage() {
  const { isAuthenticated } = useAuth();
  const [refreshNonce] = useState(0);

  const summaryPath = isAuthenticated
    ? `/insights/summary?r=${refreshNonce}`
    : `/demo/insights/summary?r=${refreshNonce}`;
  const { data: summary } = useFetch<InsightsSummary>(summaryPath);

  return (
    <div className="min-h-screen bg-crypt text-bone crypt-grain">
      <Navbar />
      <main>
        <HeroSection summary={summary} />
        <GraveyardSection refreshNonce={refreshNonce} />
        <section id="autopsy" className="mx-auto max-w-6xl border-t border-stone-border/60 px-4 py-16 sm:px-6">
          <h2 className="font-display text-3xl font-bold text-bone">Autopsy Report</h2>
          <p className="mt-2 text-sm text-stone">
            Detailed failure analytics &amp; pattern detection arriving in the next phase.
          </p>
        </section>
      </main>
      <FooterSection />
    </div>
  );
}

export function App() {
  return (
    <MotionConfig reducedMotion="user">
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<GraveyardPage />} />
            <Route path="/login" element={<GraveyardPage />} />
            <Route path="/register" element={<GraveyardPage />} />
            <Route path="/bury" element={<GraveyardPage />} />
            <Route path="/idea/:id" element={<GraveyardPage />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </MotionConfig>
  );
}

export default App;
