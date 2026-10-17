import React from 'react';
import { useAuthStore } from '../store/useAuthStore';
import { BookOpen, LogOut, User as UserIcon } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

export const Navbar: React.FC = () => {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="bg-slate-900 border-b border-slate-800 text-slate-100 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo & Tên dự án */}
          <Link to="/" className="flex items-center space-x-3 text-brand-400 hover:text-brand-300 transition-colors">
            <BookOpen className="w-7 h-7 text-brand-500" />
            <span className="font-bold text-xl tracking-wide bg-gradient-to-r from-brand-400 to-indigo-400 bg-clip-text text-transparent">
              ScriptCraft Studio
            </span>
          </Link>

          {/* User Info & Actions */}
          <div className="flex items-center space-x-4">
            {user && (
              <div className="flex items-center space-x-3 bg-slate-800/80 border border-slate-700 px-3 py-1.5 rounded-full text-sm">
                <UserIcon className="w-4 h-4 text-brand-400" />
                <span className="font-medium text-slate-200">{user.username}</span>
                <span className="bg-brand-900/60 text-brand-300 text-xs px-2 py-0.5 rounded-full border border-brand-700/50 uppercase font-semibold">
                  {user.role}
                </span>
              </div>
            )}

            <button
              onClick={handleLogout}
              className="flex items-center space-x-1.5 text-slate-400 hover:text-red-400 text-sm font-medium px-3 py-1.5 rounded-lg hover:bg-slate-800 transition-colors"
              title="Đăng xuất"
            >
              <LogOut className="w-4 h-4" />
              <span>Đăng xuất</span>
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
};
