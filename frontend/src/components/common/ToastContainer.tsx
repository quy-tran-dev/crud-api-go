import React from 'react';
import { useToastStore, type ToastType } from '../../store/useToastStore';
import { CheckCircle2, AlertCircle, Info, AlertTriangle, X } from 'lucide-react';

export const ToastContainer: React.FC = () => {
  const { toasts, removeToast } = useToastStore();

  if (toasts.length === 0) return null;

  const getToastConfig = (type: ToastType) => {
    switch (type) {
      case 'success':
        return {
          icon: <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0" />,
          borderColor: 'border-emerald-500/40',
          bgColor: 'bg-slate-900/95',
          accentColor: 'text-emerald-400',
        };
      case 'error':
        return {
          icon: <AlertCircle className="w-5 h-5 text-rose-400 flex-shrink-0" />,
          borderColor: 'border-rose-500/40',
          bgColor: 'bg-slate-900/95',
          accentColor: 'text-rose-400',
        };
      case 'warning':
        return {
          icon: <AlertTriangle className="w-5 h-5 text-amber-400 flex-shrink-0" />,
          borderColor: 'border-amber-500/40',
          bgColor: 'bg-slate-900/95',
          accentColor: 'text-amber-400',
        };
      case 'info':
      default:
        return {
          icon: <Info className="w-5 h-5 text-sky-400 flex-shrink-0" />,
          borderColor: 'border-sky-500/40',
          bgColor: 'bg-slate-900/95',
          accentColor: 'text-sky-400',
        };
    }
  };

  return (
    <div className="fixed top-5 right-5 z-[100] flex flex-col space-y-2.5 max-w-sm w-full pointer-events-none">
      {toasts.map((t) => {
        const config = getToastConfig(t.type);
        return (
          <div
            key={t.id}
            className={`pointer-events-auto flex items-start space-x-3 p-4 rounded-2xl border ${config.borderColor} ${config.bgColor} backdrop-blur-md shadow-2xl transition-all animate-in fade-in slide-in-from-top-3 duration-200`}
          >
            {config.icon}
            <div className="flex-1 text-sm font-medium text-slate-200 leading-snug">
              {t.message}
            </div>
            <button
              onClick={() => removeToast(t.id)}
              className="text-slate-400 hover:text-white p-0.5 rounded-lg hover:bg-slate-800 transition"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        );
      })}
    </div>
  );
};
