import React from 'react';
import { Dashboard } from './pages/Dashboard';

export const App: React.FC = () => {
  return (
    <div className="min-h-screen bg-redwood-bg font-sans text-redwood-text selection:bg-redwood-red/20">
      <div className="relative z-10">
        <Dashboard />
      </div>
    </div>
  );
};

export default App;
