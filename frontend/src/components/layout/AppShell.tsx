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
      className="flex min-h-screen font-sans relative bg-[#312D2A]"
      style={{
        backgroundImage: "url('/redwood-brand-bg.png')",
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundAttachment: 'fixed'
      }}
    >
      {/* Subtler frosted glass overlay allows the organic textures to shine through clearly while providing contrast for text */}
      <div className="absolute inset-0 bg-neutral-base/20 backdrop-blur-sm z-0 pointer-events-none" />

      <div className="relative z-10 flex w-full">
        <Sidebar collapsed={collapsed} setCollapsed={setCollapsed} />
        
        <div className="flex-1 flex flex-col min-w-0">
          <TopBar />
          
          <main className="flex-1 overflow-auto">
            {children}
          </main>
          
          {/* Footer info (Redwood style) */}
          <footer className="py-4 px-8 border-t border-neutral-300/30 bg-white/20 backdrop-blur-2xl text-[10px] text-text-sub flex justify-between">
            <span>&copy; 2026 Enterprise Performance Labs</span>
            <div className="flex gap-4">
              <a href="#" className="hover:underline">Legal</a>
              <a href="#" className="hover:underline">Privacy</a>
              <a href="#" className="hover:underline">Support</a>
            </div>
          </footer>
        </div>
      </div>
    </div>
  );
};
