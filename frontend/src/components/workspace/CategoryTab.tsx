import React, { useState } from 'react';
import type { SystemCategory } from '../../types';
import { api } from '../../services/api';
import { Plus, Tag, X } from 'lucide-react';

interface CategoryTabProps {
  storyID: number;
  categories: SystemCategory[];
  onRefresh: () => void;
}

export const CategoryTab: React.FC<CategoryTabProps> = ({
  storyID,
  categories,
  onRefresh,
}) => {
  const [selectedCat, setSelectedCat] = useState<SystemCategory | null>(
    categories[0] || null
  );
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [itemName, setItemName] = useState('');
  const [orderIndex, setOrderIndex] = useState(0);
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const activeCategory =
    categories.find((c) => c.id === selectedCat?.id) || categories[0] || null;

  const handleAddItem = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!activeCategory) return;
    setIsSubmitting(true);

    try {
      await api.post(`/stories/${storyID}/categories/${activeCategory.id}/items`, {
        name: itemName,
        order_index: orderIndex,
        description,
      });
      setItemName('');
      setDescription('');
      setIsModalOpen(false);
      onRefresh();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Lỗi khi thêm mục từ điển');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex h-full">
      {/* Category List */}
      <div className="w-72 border-r border-slate-800 p-4 space-y-2 overflow-y-auto">
        <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3">
          Danh Mục Hằng Số
        </h3>
        {categories.map((cat) => (
          <button
            key={cat.id}
            onClick={() => setSelectedCat(cat)}
            className={`w-full text-left p-3 rounded-xl border text-xs transition-all ${
              activeCategory?.id === cat.id
                ? 'bg-slate-800/90 border-brand-500 text-white'
                : 'bg-slate-900/60 border-slate-800 text-slate-400 hover:border-slate-700'
            }`}
          >
            <div className="font-semibold text-slate-200">{cat.name}</div>
            <div className="text-[11px] text-brand-400/80 font-mono mt-0.5">{cat.code}</div>
            <div className="text-[10px] text-slate-500 line-clamp-1 mt-1">
              {cat.items?.length || 0} mục đã định nghĩa
            </div>
          </button>
        ))}
      </div>

      {/* Items Viewer */}
      <div className="flex-1 p-6 flex flex-col overflow-y-auto">
        {activeCategory ? (
          <div>
            <div className="flex items-center justify-between pb-4 border-b border-slate-800 mb-6">
              <div>
                <h3 className="text-xl font-bold text-white flex items-center space-x-2">
                  <Tag className="w-5 h-5 text-brand-400" />
                  <span>{activeCategory.name}</span>
                </h3>
                <p className="text-xs text-slate-400 mt-0.5">
                  Mã: <code className="text-brand-300 bg-brand-950 px-1.5 py-0.5 rounded border border-brand-800/40">{activeCategory.code}</code> — {activeCategory.description}
                </p>
              </div>

              <button
                onClick={() => {
                  setOrderIndex((activeCategory.items?.length || 0) + 1);
                  setIsModalOpen(true);
                }}
                className="inline-flex items-center space-x-1.5 bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold px-3 py-2 rounded-xl transition-all shadow"
              >
                <Plus className="w-4 h-4" />
                <span>Thêm Mục Mới</span>
              </button>
            </div>

            {/* Grid of items */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {activeCategory.items?.map((item) => (
                <div
                  key={item.id}
                  className="bg-slate-900 border border-slate-800 rounded-xl p-3.5 flex items-center justify-between"
                >
                  <div>
                    <div className="font-semibold text-sm text-slate-200">
                      {item.name}
                    </div>
                    {item.description && (
                      <div className="text-xs text-slate-500 mt-0.5">
                        {item.description}
                      </div>
                    )}
                  </div>
                  <span className="text-xs bg-slate-800 text-slate-400 px-2 py-0.5 rounded-md font-mono border border-slate-700">
                    Bậc #{item.order_index}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className="text-center py-20 text-slate-500">
            Chưa có danh mục nào được chọn
          </div>
        )}
      </div>

      {/* Modal thêm mục mới */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => setIsModalOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4">
              Thêm Mục Vào Từ Điển "{activeCategory?.name}"
            </h3>
            <form onSubmit={handleAddItem} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Mục / Giá Trị *
                </label>
                <input
                  type="text"
                  required
                  value={itemName}
                  onChange={(e) => setItemName(e.target.value)}
                  placeholder="Ví dụ: Kim Đan Kỳ, Tông Môn Chính Đạo..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Thứ Tự / Cấp Bậc (Order Index)
                </label>
                <input
                  type="number"
                  value={orderIndex}
                  onChange={(e) => setOrderIndex(parseInt(e.target.value) || 0)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả / Ý Nghĩa
                </label>
                <input
                  type="text"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Ghi chú chi tiết về mục này..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>
              <div className="flex justify-end space-x-2 pt-2">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-3 py-1.5 text-xs text-slate-400 hover:text-white"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold px-4 py-2 rounded-xl transition-all"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Thêm Vào Từ Điển'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
