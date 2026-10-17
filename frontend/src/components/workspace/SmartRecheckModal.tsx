import React, { useState } from 'react';
import type { DetectedNewItem, SystemCategory } from '../../types';
import { api } from '../../services/api';
import { Sparkles, Plus, Check, X, BookText, Loader2, ArrowRight } from 'lucide-react';

interface SmartRecheckModalProps {
  isOpen: boolean;
  onClose: () => void;
  detectedItems: DetectedNewItem[];
  storyID: number;
  categories: SystemCategory[];
  onSuccess: () => void;
}

export const SmartRecheckModal: React.FC<SmartRecheckModalProps> = ({
  isOpen,
  onClose,
  detectedItems,
  storyID,
  categories,
  onSuccess,
}) => {
  const [itemsToProcess, setItemsToProcess] = useState<DetectedNewItem[]>([]);
  const [addedItems, setAddedItems] = useState<string[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);

  // Sync khi mở modal
  React.useEffect(() => {
    if (isOpen) {
      setItemsToProcess(detectedItems);
      setAddedItems([]);
    }
  }, [isOpen, detectedItems]);

  if (!isOpen || itemsToProcess.length === 0) return null;

  const handleAddSingle = async (item: DetectedNewItem) => {
    const cat = categories.find((c) => c.code === item.category_code);
    if (!cat) {
      alert(`Không tìm thấy danh mục hệ thống với mã "${item.category_code}"`);
      return;
    }

    try {
      await api.post(`/stories/${storyID}/categories/${cat.id}/items`, {
        name: item.new_value,
        order_index: (cat.items?.length || 0) + 1,
        description: 'Tự động tạo bởi tính năng Smart Recheck khi khởi tạo nhân vật',
      });

      setAddedItems((prev) => [...prev, `${item.category_code}:${item.new_value}`]);
      onSuccess();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Lỗi khi thêm mục vào từ điển');
    }
  };

  const handleAddAll = async () => {
    setIsProcessing(true);
    try {
      for (const item of itemsToProcess) {
        const cat = categories.find((c) => c.code === item.category_code);
        if (cat && !addedItems.includes(`${item.category_code}:${item.new_value}`)) {
          await api.post(`/stories/${storyID}/categories/${cat.id}/items`, {
            name: item.new_value,
            order_index: (cat.items?.length || 0) + 1,
            description: 'Tự động tạo bởi tính năng Smart Recheck',
          });
        }
      }
      onSuccess();
      onClose();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Có lỗi khi tự động thêm từ khóa');
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-amber-500/30 rounded-2xl w-full max-w-lg p-6 shadow-2xl relative">
        {/* Header */}
        <div className="flex items-start justify-between pb-4 border-b border-slate-800">
          <div className="flex items-center space-x-3">
            <div className="p-2.5 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-400">
              <Sparkles className="w-5 h-5 animate-pulse" />
            </div>
            <div>
              <h3 className="text-base font-bold text-white flex items-center space-x-2">
                <span>Smart Recheck: Phát Hiện Thuật Ngữ Mới</span>
              </h3>
              <p className="text-xs text-slate-400 mt-0.5">
                Các từ khóa vừa nhập chưa có trong Từ điển Master Data kịch bản
              </p>
            </div>
          </div>

          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* List of items */}
        <div className="mt-4 space-y-3 max-h-72 overflow-y-auto pr-1">
          {itemsToProcess.map((item, idx) => {
            const cat = categories.find((c) => c.code === item.category_code);
            const isAdded = addedItems.includes(`${item.category_code}:${item.new_value}`);

            return (
              <div
                key={idx}
                className="bg-slate-950/70 border border-slate-800 rounded-xl p-3.5 flex items-center justify-between gap-3"
              >
                <div className="space-y-1">
                  <div className="flex items-center space-x-2">
                    <span className="text-sm font-bold text-amber-300 font-mono">
                      {item.new_value}
                    </span>
                    <ArrowRight className="w-3.5 h-3.5 text-slate-600" />
                    <span className="text-xs px-2 py-0.5 rounded-full bg-slate-800 text-slate-300 border border-slate-700 flex items-center space-x-1">
                      <BookText className="w-3 h-3 text-indigo-400" />
                      <span>{cat ? cat.name : item.category_code}</span>
                    </span>
                  </div>
                  <p className="text-[11px] text-slate-400 leading-tight">
                    {item.message}
                  </p>
                </div>

                <div>
                  {isAdded ? (
                    <span className="flex items-center space-x-1 text-xs text-emerald-400 font-semibold px-2 py-1 bg-emerald-500/10 rounded-lg">
                      <Check className="w-3.5 h-3.5" />
                      <span>Đã thêm</span>
                    </span>
                  ) : (
                    <button
                      onClick={() => handleAddSingle(item)}
                      className="flex items-center space-x-1 px-2.5 py-1.5 rounded-lg bg-amber-500/20 text-amber-300 hover:bg-amber-500/30 border border-amber-500/30 text-xs font-medium transition"
                    >
                      <Plus className="w-3.5 h-3.5" />
                      <span>Thêm</span>
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>

        {/* Action footer */}
        <div className="mt-6 pt-4 border-t border-slate-800 flex items-center justify-between">
          <p className="text-[11px] text-slate-500">
            Thêm vào từ điển giúp chuẩn hóa kịch bản và gợi ý cho các nhân vật sau.
          </p>

          <div className="flex items-center space-x-2">
            <button
              onClick={onClose}
              className="px-3.5 py-2 text-xs font-medium text-slate-400 hover:text-white rounded-xl hover:bg-slate-800"
            >
              Bỏ qua
            </button>
            <button
              onClick={handleAddAll}
              disabled={isProcessing}
              className="btn-primary text-xs flex items-center space-x-1.5 bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 text-slate-950 font-bold"
            >
              {isProcessing ? (
                <>
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  <span>Đang ghi nhận...</span>
                </>
              ) : (
                <>
                  <Sparkles className="w-3.5 h-3.5" />
                  <span>Thêm Tất Cả Vào Từ Điển</span>
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
