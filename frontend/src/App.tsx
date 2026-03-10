import React from 'react';
import { Dashboard } from './pages/Dashboard';
import { AppShell } from './components/layout/AppShell';

export const App: React.FC = () => {
  return (
    <AppShell>
      <Dashboard />
    </AppShell>
  );
};

export default App;
