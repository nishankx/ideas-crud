import { AnimatePresence, MotionConfig } from 'framer-motion';
import { useCallback, useState } from 'react';
import { BrowserRouter, Navigate, Route, Routes, useLocation } from 'react-router-dom';
import { Navbar } from './components/Navbar';
import { AuthProvider, useAuth } from './hooks/useAuth';
import { useFetch } from './hooks/useFetch';
import { AuthModal } from './routes/AuthModal';
import { BuryModal } from './routes/BuryModal';
import { IdeaDetailModal } from './routes/IdeaDetailModal';
import { AutopsySection } from './sections/AutopsySection';
import { FooterSection } from './sections/FooterSection';
import { GraveyardSection } from './sections/GraveyardSection';
import { HeroSection } from './sections/HeroSection';
import type { InsightsSummary } from './types/api';

function GraveyardShell() {
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  const [refreshNonce, setRefreshNonce] = useState(0);

  const bumpRefresh = useCallback(() => setRefreshNonce((n) => n + 1), []);

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
        <AutopsySection refreshNonce={refreshNonce} summary={summary} />
      </main>
      <FooterSection />

      <AnimatePresence mode="wait">
        <Routes location={location} key={location.pathname}>
          <Route path="/" element={null} />
          <Route path="/login" element={<AuthModal mode="login" />} />
          <Route path="/register" element={<AuthModal mode="register" />} />
          <Route path="/bury" element={<BuryModal onBuried={bumpRefresh} />} />
          <Route path="/idea/:id" element={<IdeaDetailModal onMutated={bumpRefresh} />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AnimatePresence>
    </div>
  );
}

export function App() {
  return (
    <MotionConfig reducedMotion="user">
      <AuthProvider>
        <BrowserRouter>
          <GraveyardShell />
        </BrowserRouter>
      </AuthProvider>
    </MotionConfig>
  );
}

export default App;
