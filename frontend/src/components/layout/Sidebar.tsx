import React, { useState } from 'react';
import { 
  LayoutDashboard, 
  Activity, 
  Settings, 
  ChevronLeft, 
  ChevronRight, 
  Database,
  ShieldCheck,
  Zap
} from 'lucide-react';
import { cn } from '../../utils/cn';

interface SidebarProps {
  collapsed: boolean;
  setCollapsed: (val: boolean) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ collapsed, setCollapsed }) => {
  const [activeItem, setActiveItem] = useState('dashboard');

  const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { id: 'latency', label: 'Latency Map', icon: Activity },
    { id: 'security', label: 'Security Portal', icon: ShieldCheck },
    { id: 'data', label: 'Data Streams', icon: Database },
    { id: 'perf', label: 'Benchmark Lab', icon: Zap },
    { id: 'settings', label: 'System Config', icon: Settings },
  ];

  return (
    <aside 
      className={cn(
        "h-screen sticky top-0 bg-neutral-base border-r border-neutral-200 transition-all duration-300 z-50 flex flex-col",
        collapsed ? "w-16" : "w-64"
      )}
    >
      {/* Sidebar Header */}
      <div className="h-16 flex items-center px-4 border-b border-neutral-200">
        {!collapsed && (
          <span className="text-primary-red font-black tracking-tighter text-xl ml-2">RW-LAB</span>
        )}
        <button 
          onClick={() => setCollapsed(!collapsed)}
          className="ml-auto p-2 hover:bg-neutral-100 rounded-md transition-colors"
        >
          {collapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
        </button>
      </div>

      {/* Nav Items */}
      <nav className="flex-1 py-4 overflow-y-auto overflow-x-hidden">
        {navItems.map((item) => (
          <div 
            key={item.id}
            onClick={() => setActiveItem(item.id)}
            className={cn(
              "rw-nav-item mx-2 mb-1",
              activeItem === item.id && "rw-nav-item-active"
            )}
          >
            <item.icon size={20} className={cn(collapsed ? "mx-auto" : "min-w-[20px]")} />
            {!collapsed && <span className="truncate">{item.label}</span>}
          </div>
        ))}
      </nav>

      {/* Footer / Identity */}
      <div className="p-4 border-t border-neutral-200">
         <div className={cn("flex items-center gap-3", collapsed ? "justify-center" : "")}>
            <div className="w-8 h-8 rounded-full bg-neutral-200 flex items-center justify-center text-[10px] font-bold">AS</div>
            {!collapsed && (
              <div className="flex flex-col">
                <span className="text-[12px] font-bold truncate">Andy S.</span>
                <span className="text-[10px] text-text-sub truncate">Performance Lead</span>
              </div>
            )}
         </div>
      </div>
    </aside>
  );
};
