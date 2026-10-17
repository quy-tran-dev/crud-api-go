import React from 'react';
import {
  BookText,
  MapPin,
  Shield,
  Users,
  Film,
  Sword,
  Folder,
  FolderOpen,
} from 'lucide-react';

export type WorkspaceTab =
  | 'categories'
  | 'locations'
  | 'factions'
  | 'characters'
  | 'plots'
  | 'items';

interface WorkspaceSidebarProps {
  activeTab: WorkspaceTab;
  onTabChange: (tab: WorkspaceTab) => void;
  counts: {
    categories: number;
    locations: number;
    factions: number;
    characters: number;
    plots: number;
    items: number;
  };
}

export const WorkspaceSidebar: React.FC<WorkspaceSidebarProps> = ({
  activeTab,
  onTabChange,
  counts,
}) => {
  const folders: { id: WorkspaceTab; label: string; icon: any; count: number; color: string }[] = [
    {
      id: 'characters',
      label: 'Nhân Vật',
      icon: Users,
      count: counts.characters,
      color: 'text-amber-400',
    },
    {
      id: 'locations',
      label: 'Bản Đồ & Địa Danh',
      icon: MapPin,
      count: counts.locations,
      color: 'text-emerald-400',
    },
    {
      id: 'factions',
      label: 'Thế Lực & Gia Tộc',
      icon: Shield,
      count: counts.factions,
      color: 'text-blue-400',
    },
    {
      id: 'plots',
      label: 'Tình Tiết (Timeline)',
      icon: Film,
      count: counts.plots,
      color: 'text-purple-400',
    },
    {
      id: 'items',
      label: 'Vật Phẩm & Phẩm Cấp',
      icon: Sword,
      count: counts.items,
      color: 'text-rose-400',
    },
    {
      id: 'categories',
      label: 'Từ Điển Hằng Số',
      icon: BookText,
      count: counts.categories,
      color: 'text-indigo-400',
    },
  ];

  return (
    <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col flex-shrink-0 select-none">
      <div className="p-4 border-b border-slate-800">
        <h2 className="text-xs font-bold text-slate-400 uppercase tracking-wider flex items-center space-x-2">
          <Folder className="w-4 h-4 text-brand-400" />
          <span>Cây Thư Mục Kịch Bản</span>
        </h2>
      </div>

      <nav className="flex-1 p-2 space-y-1 overflow-y-auto">
        {folders.map((f) => {
          const isActive = activeTab === f.id;
          const Icon = f.icon;

          return (
            <button
              key={f.id}
              onClick={() => onTabChange(f.id)}
              className={`w-full flex items-center justify-between px-3 py-2.5 rounded-xl text-xs font-medium transition-all ${
                isActive
                  ? 'bg-brand-600/20 text-white border border-brand-500/40 shadow-sm'
                  : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
              }`}
            >
              <div className="flex items-center space-x-2.5 truncate">
                {isActive ? (
                  <FolderOpen className={`w-4 h-4 ${f.color}`} />
                ) : (
                  <Icon className={`w-4 h-4 ${f.color} opacity-80`} />
                )}
                <span className="truncate">{f.label}</span>
              </div>

              <span
                className={`text-[11px] px-2 py-0.5 rounded-full font-semibold ${
                  isActive
                    ? 'bg-brand-500/30 text-brand-300'
                    : 'bg-slate-800 text-slate-500'
                }`}
              >
                {f.count}
              </span>
            </button>
          );
        })}
      </nav>
    </aside>
  );
};
