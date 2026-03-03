import React from 'react';
import { Dashboard } from './pages/Dashboard';

export const App: React.FC = () => {
  return (
    <div className="min-h-screen bg-zinc-950 font-sans text-zinc-100 selection:bg-indigo-500/30">
      <div className="absolute inset-0 bg-[url('https://res.cloudinary.com/dhdq1m04t/image/upload/v1689617267/grid_b2k10i.svg')] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))] opacity-5 pointer-events-none" />
      <div className="relative z-10">
        <Dashboard />
      </div>
    </div>
  );
};

export default App;
