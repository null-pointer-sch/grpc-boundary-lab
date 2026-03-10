import React, { useState } from 'react';
import { Sidebar } from './Sidebar';
import { TopBar } from './TopBar';

interface AppShellProps {
  children: React.ReactNode;
}

export const AppShell: React.FC<AppShellProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div 
      className="flex min-h-screen font-sans relative"
      style={{
        backgroundColor: '#312D2A',
        backgroundImage: "url('/redwood-brand-bg.png')",
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundAttachment: 'fixed'
      }}
    >
      <div className="relative z-10 flex w-full p-6 gap-6 max-w-[1920px] mx-auto">
        <div className="rounded-xl overflow-hidden shadow-2xl flex-shrink-0">
          <Sidebar collapsed={collapsed} setCollapsed={setCollapsed} />
        </div>
        
        <div className="flex-1 flex flex-col min-w-0 bg-[#F4EBE1] rounded-xl overflow-hidden shadow-2xl">
          <TopBar />
          
          <main className="flex-1 overflow-auto p-8">
            {children}
          </main>
          
          <footer className="py-4 px-8 border-t border-neutral-300/30 text-[10px] text-text-sub flex justify-between">
            <span>&copy; 2026 Enterprise Performance Labs</span>
            <div className="flex gap-4">
              <button className="hover:underline">Legal</button>
              <button className="hover:underline">Privacy</button>
              <button className="hover:underline">Support</button>
            </div>
          </footer>
        </div>
      </div>
    </div>
  );
};
