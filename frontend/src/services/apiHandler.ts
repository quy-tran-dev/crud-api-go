import { toast } from '../store/useToastStore';
import { confirmModal } from '../store/useConfirmStore';

/**
 * Trích xuất thông báo lỗi chuẩn từ response của Go Backend hoặc fallback message
 */
export const getErrorMessage = (
  err: any,
  fallbackMessage: string = 'Đã có lỗi xảy ra. Vui lòng thử lại.'
): string => {
  return (
    err?.response?.data?.message ||
    err?.response?.data?.error ||
    err?.message ||
    fallbackMessage
  );
};

export interface HandleFormSubmitOptions<T> {
  action: () => Promise<T>;
  setSubmitting?: (isSubmitting: boolean) => void;
  onSuccess?: (result: T) => void;
  onError?: (errorMessage: string) => void;
  successMessage?: string;
  fallbackErrorMessage?: string;
  showErrorToast?: boolean;
}

/**
 * Xử lý vòng đời submit form chuẩn hóa (bật loading, gọi API, bắt lỗi, toast thông báo, tắt loading)
 */
export const handleFormSubmit = async <T>({
  action,
  setSubmitting,
  onSuccess,
  onError,
  successMessage,
  fallbackErrorMessage = 'Lỗi khi lưu dữ liệu',
  showErrorToast = true,
}: HandleFormSubmitOptions<T>): Promise<T | null> => {
  if (setSubmitting) setSubmitting(true);
  try {
    const result = await action();
    if (successMessage) {
      toast.success(successMessage);
    }
    if (onSuccess) onSuccess(result);
    return result;
  } catch (err: any) {
    const msg = getErrorMessage(err, fallbackErrorMessage);
    if (onError) {
      onError(msg);
    }
    if (showErrorToast) {
      toast.error(msg);
    }
    return null;
  } finally {
    if (setSubmitting) setSubmitting(false);
  }
};

export interface ConfirmAndDeleteOptions {
  title?: string;
  confirmMessage: string;
  confirmText?: string;
  cancelText?: string;
  action: () => Promise<any>;
  onSuccess?: () => void;
  onError?: (errorMessage: string) => void;
  successMessage?: string;
  fallbackErrorMessage?: string;
}

/**
 * Hiển thị Popup xác nhận hiện đại và thực hiện hành động xóa, tự động bắn Toast thông báo
 */
export const confirmAndDelete = async ({
  title = 'Xác nhận xóa',
  confirmMessage,
  confirmText = 'Xác nhận xóa',
  cancelText = 'Hủy bỏ',
  action,
  onSuccess,
  onError,
  successMessage = 'Đã xóa dữ liệu thành công!',
  fallbackErrorMessage = 'Lỗi khi xóa dữ liệu',
}: ConfirmAndDeleteOptions): Promise<boolean> => {
  const confirmed = await confirmModal({
    title,
    message: confirmMessage,
    confirmText,
    cancelText,
    variant: 'danger',
  });

  if (!confirmed) return false;

  try {
    await action();
    if (successMessage) {
      toast.success(successMessage);
    }
    if (onSuccess) onSuccess();
    return true;
  } catch (err: any) {
    const msg = getErrorMessage(err, fallbackErrorMessage);
    if (onError) {
      onError(msg);
    } else {
      toast.error(msg);
    }
    return false;
  }
};
