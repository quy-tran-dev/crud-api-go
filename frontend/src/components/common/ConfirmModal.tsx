import React from 'react';
import { useConfirmStore } from '../../store/useConfirmStore';
import { AlertTriangle, AlertCircle, Info, X } from 'lucide-react';

export const ConfirmModal: React.FC = () => {
  const { isOpen, options, handleConfirm, handleCancel } = useConfirmStore();

  if (!isOpen || !options) return null;

  const getVariantStyles = () => {
    switch (options.variant) {
      case 'warning':
        return {
          icon: <AlertCircle className="w-6 h-6 text-amber-400" />,
          iconBg: 'bg-amber-500/10 border-amber-500/20',
          confirmBtn: 'bg-amber-600 hover:bg-amber-500 text-white',
        };
      case 'info':
        return {
          icon: <Info className="w-6 h-6 text-sky-400" />,
          iconBg: 'bg-sky-500/10 border-sky-500/20',
          confirmBtn: 'bg-brand-600 hover:bg-brand-500 text-white',
        };
      case 'danger':
      default:
        return {
          icon: <AlertTriangle className="w-6 h-6 text-rose-400" />,
          iconBg: 'bg-rose-500/10 border-rose-500/20',
          confirmBtn: 'btn-danger',
        };
    }
  };

  const styles = getVariantStyles();

  return (
    <div className="fixed inset-0 z-[110] bg-black/75 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl relative animate-in zoom-in-95 duration-200">
        <button
          onClick={handleCancel}
          className="absolute top-4 right-4 text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="flex items-start space-x-4">
          <div className={`p-3 rounded-2xl border ${styles.iconBg} flex-shrink-0`}>
            {styles.icon}
          </div>

          <div className="flex-1 pt-0.5">
            <h3 className="text-base font-bold text-white mb-1.5">
              {options.title}
            </h3>
            <p className="text-sm text-slate-300 leading-relaxed whitespace-pre-wrap">
              {options.message}
            </p>
          </div>
        </div>

        <div className="flex items-center justify-end space-x-3 mt-6 pt-4 border-t border-slate-800/80">
          <button
            type="button"
            onClick={handleCancel}
            className="btn-secondary"
          >
            {options.cancelText}
          </button>
          <button
            type="button"
            onClick={handleConfirm}
            className={styles.confirmBtn}
          >
            {options.confirmText}
          </button>
        </div>
      </div>
    </div>
  );
};
