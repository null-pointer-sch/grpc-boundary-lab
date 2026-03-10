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
    <header className="h-16 bg-[#312D2A] border-b border-neutral-300/20 px-6 flex items-center justify-between sticky top-0 z-40">
      {/* Left side: Breadcrumb / Identity */}
      <div className="flex items-center gap-4">
        <div className="flex items-center text-[12px] font-medium text-white/70 gap-2">
          <span className="hover:text-white cursor-pointer transition-colors">Enterprise Home</span>
          <ChevronRight size={14} />
          <span className="text-white font-semibold">Boundary Lab</span>
        </div>
      </div>

      {/* Middle: Search */}
      <div className="hidden md:flex flex-1 max-w-md mx-8">
        <div className="relative w-full">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-white/50" size={16} />
          <input 
            type="text" 
            placeholder="Search performance metrics..." 
            className="w-full h-9 bg-white/10 border border-white/20 rounded-md pl-10 pr-4 text-sm text-white placeholder:text-white/50 focus:outline-none focus:ring-2 focus:ring-primary-red focus:border-transparent transition-all"
          />
        </div>
      </div>

      {/* Right side: Actions */}
      <div className="flex items-center gap-5 text-white/80">
        <button className="p-2 hover:bg-white/10 hover:text-white rounded-full transition-colors relative">
          <Bell size={20} />
          <span className="absolute top-2 right-2 w-2 h-2 bg-primary-red rounded-full border-2 border-[#312D2A]"></span>
        </button>
        <button className="p-2 hover:bg-white/10 hover:text-white rounded-full transition-colors">
          <HelpCircle size={20} />
        </button>
        <div className="h-8 w-px bg-white/20 mx-1"></div>
        <button className="flex items-center gap-2 px-2 py-1 hover:bg-white/10 hover:text-white rounded-md transition-colors text-white">
          <span className="text-sm font-bold">AS</span>
          <UserCircle size={24} className="text-white/50" />
        </button>
      </div>
    </header>
  );
};
