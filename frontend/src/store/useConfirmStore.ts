import { create } from 'zustand';

export interface ConfirmOptions {
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  variant?: 'danger' | 'warning' | 'info';
}

interface ConfirmState {
  isOpen: boolean;
  options: ConfirmOptions | null;
  resolveCallback: ((confirmed: boolean) => void) | null;
  openConfirm: (options: ConfirmOptions) => Promise<boolean>;
  handleConfirm: () => void;
  handleCancel: () => void;
}

export const useConfirmStore = create<ConfirmState>((set, get) => ({
  isOpen: false,
  options: null,
  resolveCallback: null,

  openConfirm: (options: ConfirmOptions) => {
    return new Promise<boolean>((resolve) => {
      set({
        isOpen: true,
        options: {
          title: options.title || 'Xác nhận thao tác',
          confirmText: options.confirmText || 'Xác nhận',
          cancelText: options.cancelText || 'Hủy bỏ',
          variant: options.variant || 'danger',
          message: options.message,
        },
        resolveCallback: resolve,
      });
    });
  },

  handleConfirm: () => {
    const { resolveCallback } = get();
    if (resolveCallback) resolveCallback(true);
    set({ isOpen: false, options: null, resolveCallback: null });
  },

  handleCancel: () => {
    const { resolveCallback } = get();
    if (resolveCallback) resolveCallback(false);
    set({ isOpen: false, options: null, resolveCallback: null });
  },
}));

// Helper tiện ích để gọi confirm modal dạng async/await ở bất cứ đâu
export const confirmModal = (options: ConfirmOptions | string): Promise<boolean> => {
  const opts = typeof options === 'string' ? { message: options } : options;
  return useConfirmStore.getState().openConfirm(opts);
};
