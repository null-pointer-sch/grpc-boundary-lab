import React from 'react';
import { 
  Bell, 
  Search, 
  HelpCircle,
  ChevronRight,
  UserCircle
} from 'lucide-react';

export const TopBar: React.FC = () => {
  return (
    <header className="h-16 bg-white border-b border-neutral-200 px-6 flex items-center justify-between sticky top-0 z-40">
      {/* Left side: Breadcrumb / Identity */}
      <div className="flex items-center gap-4">
        <div className="flex items-center text-[12px] font-medium text-text-sub gap-2">
          <span>Enterprise Home</span>
          <ChevronRight size={14} />
          <span className="text-text-main font-semibold">Boundary Lab</span>
        </div>
      </div>

      {/* Middle: Search */}
      <div className="hidden md:flex flex-1 max-w-md mx-8">
        <div className="relative w-full">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-text-sub" size={16} />
          <input 
            type="text" 
            placeholder="Search performance metrics..." 
            className="w-full h-9 bg-neutral-base border border-neutral-200 rounded-md pl-10 pr-4 text-sm focus:outline-none focus:ring-2 focus:ring-primary-red/20 focus:border-primary-red transition-all"
          />
        </div>
      </div>

      {/* Right side: Actions */}
      <div className="flex items-center gap-5 text-text-sub">
        <button className="p-2 hover:bg-neutral-100 rounded-full transition-colors relative">
          <Bell size={20} />
          <span className="absolute top-2 right-2 w-2 h-2 bg-primary-red rounded-full border-2 border-white"></span>
        </button>
        <button className="p-2 hover:bg-neutral-100 rounded-full transition-colors">
          <HelpCircle size={20} />
        </button>
        <div className="h-8 w-px bg-neutral-200 mx-1"></div>
        <button className="flex items-center gap-2 px-2 py-1 hover:bg-neutral-100 rounded-md transition-colors text-text-main">
          <span className="text-sm font-bold">AS</span>
          <UserCircle size={24} className="text-neutral-300" />
        </button>
      </div>
    </header>
  );
};
