import React, { useState } from 'react';
import { Sidebar } from './Sidebar';
import { TopBar } from './TopBar';

interface AppShellProps {
  children: React.ReactNode;
}

export const AppShell: React.FC<AppShellProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex min-h-screen bg-neutral-base font-sans">
      <Sidebar collapsed={collapsed} setCollapsed={setCollapsed} />
      
      <div className="flex-1 flex flex-col min-w-0">
        <TopBar />
        
        <main className="flex-1 overflow-auto">
          {children}
        </main>
        
        {/* Footer info (Redwood style) */}
        <footer className="py-4 px-8 border-t border-neutral-200 bg-white text-[10px] text-text-sub flex justify-between">
          <span>&copy; 2026 Enterprise Performance Labs</span>
          <div className="flex gap-4">
            <a href="#" className="hover:underline">Legal</a>
            <a href="#" className="hover:underline">Privacy</a>
            <a href="#" className="hover:underline">Support</a>
          </div>
        </footer>
      </div>
    </div>
  );
};
