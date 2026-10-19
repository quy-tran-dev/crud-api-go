import React, { useState } from 'react';
import type { SystemCategory, SystemCategoryItem } from '../../types';
import { api, ENDPOINTS, handleFormSubmit, confirmAndDelete } from '../../services';
import {
  Plus,
  Tag,
  X,
  Edit2,
  Trash2,
  BookmarkPlus,
  Share2,
  Check,
  FolderPlus,
} from 'lucide-react';

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
  const [selectedCatID, setSelectedCatID] = useState<number | null>(
    categories[0]?.id || null
  );

  // Modal states - Category
  const [isCreateCatOpen, setIsCreateCatOpen] = useState(false);
  const [isEditCatOpen, setIsEditCatOpen] = useState(false);
  const [catName, setCatName] = useState('');
  const [catCode, setCatCode] = useState('');
  const [catDesc, setCatDesc] = useState('');

  // Modal states - Category Item
  const [isCreateItemOpen, setIsCreateItemOpen] = useState(false);
  const [isEditItemOpen, setIsEditItemOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<SystemCategoryItem | null>(null);
  const [itemName, setItemName] = useState('');
  const [orderIndex, setOrderIndex] = useState(0);
  const [itemDesc, setItemDesc] = useState('');

  // Modal states - Save Framework Template
  const [isSaveTemplateOpen, setIsSaveTemplateOpen] = useState(false);
  const [tmplName, setTmplName] = useState('');
  const [tmplDesc, setTmplDesc] = useState('');
  const [tmplIsPublic, setTmplIsPublic] = useState(false);
  const [tmplSuccessMsg, setTmplSuccessMsg] = useState('');

  const [isSubmitting, setIsSubmitting] = useState(false);

  const activeCategory =
    categories.find((c) => c.id === selectedCatID) || categories[0] || null;

  // --- Category Handlers ---
  const handleCreateCategory = async (e: React.SubmitEvent) => {
    e.preventDefault();
    const code = catCode.trim() || catName.trim().toLowerCase().replace(/\s+/g, '_');
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.CATEGORIES.BASE(storyID), {
          name: catName,
          code,
          description: catDesc,
        }),
      onSuccess: (res) => {
        setCatName('');
        setCatCode('');
        setCatDesc('');
        setIsCreateCatOpen(false);
        onRefresh();
        if (res.data?.data?.id) {
          setSelectedCatID(res.data.data.id);
        }
      },
      fallbackErrorMessage: 'Lỗi khi tạo danh mục mới',
      successMessage: 'Đã tạo danh mục mới thành công!',
    });
  };

  const handleUpdateCategory = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!activeCategory) return;
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.put(ENDPOINTS.CATEGORIES.DETAIL(storyID, activeCategory.id), {
          name: catName,
          description: catDesc,
        }),
      onSuccess: () => {
        setIsEditCatOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi cập nhật danh mục',
      successMessage: 'Đã cập nhật danh mục thành công!',
    });
  };

  const handleDeleteCategory = async (cat: SystemCategory) => {
    await confirmAndDelete({
      title: 'Xóa Danh Mục',
      confirmMessage: `Bạn có chắc chắn muốn xóa danh mục "${cat.name}" cùng tất cả các mục bên trong?`,
      action: () => api.delete(ENDPOINTS.CATEGORIES.DETAIL(storyID, cat.id)),
      onSuccess: () => {
        onRefresh();
        const remain = categories.filter((c) => c.id !== cat.id);
        if (remain.length > 0) {
          setSelectedCatID(remain[0].id);
        } else {
          setSelectedCatID(null);
        }
      },
      fallbackErrorMessage: 'Lỗi khi xóa danh mục',
      successMessage: `Đã xóa danh mục "${cat.name}" thành công!`,
    });
  };

  // --- Category Item Handlers ---
  const handleCreateItem = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!activeCategory) return;
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.CATEGORIES.ITEMS.BASE(storyID, activeCategory.id), {
          name: itemName,
          order_index: Number(orderIndex),
          description: itemDesc,
        }),
      onSuccess: () => {
        setItemName('');
        setItemDesc('');
        setIsCreateItemOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi thêm mục từ điển',
      successMessage: 'Đã thêm mục từ điển thành công!',
    });
  };

  const handleUpdateItem = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!activeCategory || !editingItem) return;
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.put(
          ENDPOINTS.CATEGORIES.ITEMS.DETAIL(
            storyID,
            activeCategory.id,
            editingItem.id
          ),
          {
            name: itemName,
            order_index: Number(orderIndex),
            description: itemDesc,
          }
        ),
      onSuccess: () => {
        setIsEditItemOpen(false);
        setEditingItem(null);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi cập nhật mục từ điển',
      successMessage: 'Đã cập nhật mục từ điển thành công!',
    });
  };

  const handleDeleteItem = async (itemID: number, name: string) => {
    if (!activeCategory) return;
    await confirmAndDelete({
      title: 'Xóa Mục Từ Điển',
      confirmMessage: `Bạn có chắc muốn xóa mục "${name}" khỏi từ điển?`,
      action: () =>
        api.delete(
          ENDPOINTS.CATEGORIES.ITEMS.DETAIL(storyID, activeCategory.id, itemID)
        ),
      onSuccess: () => {
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi xóa mục từ điển',
      successMessage: `Đã xóa mục "${name}" thành công!`,
    });
  };

  // --- Save Framework Template Handler ---
  const handleSaveFramework = async (e: React.SubmitEvent) => {
    e.preventDefault();
    setTmplSuccessMsg('');
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.STORIES.SAVE_FRAMEWORK(storyID), {
          name: tmplName,
          description: tmplDesc,
          is_public: tmplIsPublic,
        }),
      onSuccess: (res) => {
        setTmplSuccessMsg(
          res.data?.message ||
            'Đã lưu bộ khung thành công! Bạn có thể tái sử dụng hoặc chia sẻ với cộng đồng.'
        );
        setTimeout(() => {
          setIsSaveTemplateOpen(false);
          setTmplSuccessMsg('');
          setTmplName('');
          setTmplDesc('');
        }, 1500);
      },
      fallbackErrorMessage: 'Lỗi khi lưu bộ khung',
      successMessage: 'Đã lưu bộ khung thành công!',
    });
  };

  // Sort items by order_index
  const sortedItems = [...(activeCategory?.items || [])].sort(
    (a, b) => a.order_index - b.order_index
  );

  return (
    <div className="flex h-full overflow-hidden">
      {/* Category List Sidebar */}
      <div className="w-72 border-r border-slate-800 p-4 flex flex-col justify-between overflow-y-auto flex-shrink-0 bg-slate-900/40">
        <div>
          <div className="flex items-center justify-between mb-3 pb-2 border-b border-slate-800">
            <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider">
              Danh Mục Hằng Số
            </h3>
            <button
              onClick={() => {
                setCatName('');
                setCatCode('');
                setCatDesc('');
                setIsCreateCatOpen(true);
              }}
              className="p-1 rounded-lg text-brand-400 hover:text-white hover:bg-brand-600/20 transition flex items-center space-x-1 text-xs"
              title="Thêm Danh Mục Mới"
            >
              <FolderPlus className="w-4 h-4" />
              <span className="font-semibold text-[11px]">Thêm</span>
            </button>
          </div>

          <div className="space-y-1.5">
            {categories.map((cat) => {
              const isActive = activeCategory?.id === cat.id;
              return (
                <div
                  key={cat.id}
                  onClick={() => setSelectedCatID(cat.id)}
                  className={`w-full text-left p-3 rounded-xl border text-xs transition-all cursor-pointer group flex items-start justify-between ${
                    isActive
                      ? 'bg-brand-950/40 border-brand-500/60 text-white shadow-sm'
                      : 'bg-slate-900/60 border-slate-800 text-slate-400 hover:border-slate-700 hover:text-slate-200'
                  }`}
                >
                  <div className="truncate flex-1 mr-2">
                    <div className="font-semibold text-slate-200 truncate">
                      {cat.name}
                    </div>
                    <div className="text-[11px] text-brand-400/80 font-mono mt-0.5 truncate">
                      {cat.code}
                    </div>
                    <div className="text-[10px] text-slate-500 line-clamp-1 mt-1">
                      {cat.items?.length || 0} mục đã định nghĩa
                    </div>
                  </div>

                  {/* Actions on category */}
                  <div className="opacity-0 group-hover:opacity-100 flex items-center space-x-1 transition-opacity">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setCatName(cat.name);
                        setCatDesc(cat.description || '');
                        setSelectedCatID(cat.id);
                        setIsEditCatOpen(true);
                      }}
                      className="p-1 rounded hover:bg-slate-700 text-slate-400 hover:text-white"
                      title="Sửa danh mục"
                    >
                      <Edit2 className="w-3.5 h-3.5" />
                    </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDeleteCategory(cat);
                      }}
                      className="p-1 rounded hover:bg-rose-950/60 text-slate-400 hover:text-rose-400"
                      title="Xóa danh mục"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Action Button: Save entire world categories as Template */}
        <div className="pt-4 border-t border-slate-800/80">
          <button
            onClick={() => {
              setTmplName('');
              setTmplDesc('');
              setTmplIsPublic(false);
              setTmplSuccessMsg('');
              setIsSaveTemplateOpen(true);
            }}
            className="w-full flex items-center justify-center space-x-1.5 py-2.5 px-3 rounded-xl bg-gradient-to-r from-indigo-600/30 to-purple-600/30 hover:from-indigo-600/50 hover:to-purple-600/50 border border-indigo-500/30 text-indigo-300 text-xs font-semibold transition"
          >
            <BookmarkPlus className="w-4 h-4 text-indigo-400" />
            <span>Lưu Thành Bộ Khung Mẫu</span>
          </button>
        </div>
      </div>

      {/* Items Viewer */}
      <div className="flex-1 p-6 flex flex-col overflow-y-auto bg-slate-950">
        {activeCategory ? (
          <div>
            {/* Header */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-800 mb-6 gap-3">
              <div>
                <div className="flex items-center space-x-2.5">
                  <Tag className="w-5 h-5 text-brand-400" />
                  <h3 className="text-xl font-bold text-white">
                    {activeCategory.name}
                  </h3>
                  <button
                    onClick={() => {
                      setCatName(activeCategory.name);
                      setCatDesc(activeCategory.description || '');
                      setIsEditCatOpen(true);
                    }}
                    className="p-1 text-slate-400 hover:text-white rounded hover:bg-slate-800 transition"
                    title="Chỉnh sửa danh mục này"
                  >
                    <Edit2 className="w-4 h-4" />
                  </button>
                </div>
                <p className="text-xs text-slate-400 mt-1">
                  Mã:{' '}
                  <code className="text-brand-300 bg-brand-950 px-1.5 py-0.5 rounded border border-brand-800/40">
                    {activeCategory.code}
                  </code>{' '}
                  — {activeCategory.description || 'Không có mô tả'}
                </p>
              </div>

              <div className="flex items-center space-x-2">
                <button
                  onClick={() => {
                    setItemName('');
                    setItemDesc('');
                    setOrderIndex(sortedItems.length + 1);
                    setIsCreateItemOpen(true);
                  }}
                  className="btn-primary text-xs flex items-center space-x-1.5 py-2 px-3"
                >
                  <Plus className="w-4 h-4" />
                  <span>Thêm Cấp Bậc / Mục Mới</span>
                </button>
              </div>
            </div>

            {/* Grid of items */}
            {sortedItems.length === 0 ? (
              <div className="bg-slate-900 border border-dashed border-slate-800 rounded-2xl p-12 text-center text-slate-400">
                <Tag className="w-10 h-10 mx-auto mb-3 text-slate-600" />
                <p className="font-semibold text-slate-300">
                  Chưa có mục hay cấp bậc nào trong danh mục này
                </p>
                <p className="text-xs mt-1 text-slate-500">
                  Hãy thêm cấp bậc đầu tiên (ví dụ: Luyện Khí Kỳ, Kim Đan Kỳ...)
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                {sortedItems.map((item) => (
                  <div
                    key={item.id}
                    className="bg-slate-900 border border-slate-800 rounded-xl p-3.5 flex items-start justify-between group hover:border-slate-700 transition"
                  >
                    <div className="flex-1 mr-2">
                      <div className="flex items-center space-x-2">
                        <span className="text-xs bg-slate-800 text-brand-400 font-mono font-bold px-2 py-0.5 rounded border border-slate-700/60">
                          #{item.order_index}
                        </span>
                        <div className="font-semibold text-sm text-slate-200">
                          {item.name}
                        </div>
                      </div>
                      {item.description && (
                        <div className="text-xs text-slate-400 mt-2 leading-relaxed bg-slate-950/40 p-2 rounded-lg border border-slate-800/60">
                          {item.description}
                        </div>
                      )}
                    </div>

                    <div className="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        onClick={() => {
                          setEditingItem(item);
                          setItemName(item.name);
                          setOrderIndex(item.order_index);
                          setItemDesc(item.description || '');
                          setIsEditItemOpen(true);
                        }}
                        className="p-1 rounded hover:bg-slate-800 text-slate-400 hover:text-white"
                        title="Chỉnh sửa cấp bậc"
                      >
                        <Edit2 className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => handleDeleteItem(item.id, item.name)}
                        className="p-1 rounded hover:bg-rose-950/60 text-slate-400 hover:text-rose-400"
                        title="Xóa cấp bậc"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        ) : (
          <div className="text-center py-20 text-slate-500">
            Chưa có danh mục nào được chọn
          </div>
        )}
      </div>

      {/* --- Modal: Tạo Danh Mục Mới --- */}
      {isCreateCatOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => setIsCreateCatOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4 flex items-center space-x-2">
              <FolderPlus className="w-5 h-5 text-brand-400" />
              <span>Tạo Danh Mục Hằng Số Mới</span>
            </h3>
            <form onSubmit={handleCreateCategory} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Danh Mục *
                </label>
                <input
                  type="text"
                  required
                  value={catName}
                  onChange={(e) => setCatName(e.target.value)}
                  placeholder="Ví dụ: Hệ Thống Linh Căn, Chủng Tộc..."
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mã Code (Tùy chọn, tự động tạo nếu bỏ trống)
                </label>
                <input
                  type="text"
                  value={catCode}
                  onChange={(e) => setCatCode(e.target.value)}
                  placeholder="Ví dụ: element_type, race..."
                  className="input-field font-mono text-xs"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả Danh Mục
                </label>
                <textarea
                  rows={2}
                  value={catDesc}
                  onChange={(e) => setCatDesc(e.target.value)}
                  placeholder="Mục đích sử dụng của danh mục này trong thế giới..."
                  className="input-field text-xs"
                />
              </div>
              <div className="flex justify-end space-x-2 pt-2 border-t border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsCreateCatOpen(false)}
                  className="px-3 py-1.5 text-xs text-slate-400 hover:text-white"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="btn-primary text-xs"
                >
                  {isSubmitting ? 'Đang tạo...' : 'Tạo Danh Mục'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- Modal: Sửa Danh Mục --- */}
      {isEditCatOpen && activeCategory && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => setIsEditCatOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4 flex items-center space-x-2">
              <Edit2 className="w-5 h-5 text-brand-400" />
              <span>Chỉnh Sửa Danh Mục: {activeCategory.code}</span>
            </h3>
            <form onSubmit={handleUpdateCategory} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Danh Mục *
                </label>
                <input
                  type="text"
                  required
                  value={catName}
                  onChange={(e) => setCatName(e.target.value)}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả Danh Mục
                </label>
                <textarea
                  rows={2}
                  value={catDesc}
                  onChange={(e) => setCatDesc(e.target.value)}
                  className="input-field text-xs"
                />
              </div>
              <div className="flex justify-end space-x-2 pt-2 border-t border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsEditCatOpen(false)}
                  className="px-3 py-1.5 text-xs text-slate-400 hover:text-white"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="btn-primary text-xs"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Lưu Thay Đổi'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- Modal: Thêm Cấp Bậc / Mục Mới --- */}
      {isCreateItemOpen && activeCategory && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => setIsCreateItemOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4 flex items-center space-x-2">
              <Plus className="w-5 h-5 text-brand-400" />
              <span>Thêm Vào Danh Mục "{activeCategory.name}"</span>
            </h3>
            <form onSubmit={handleCreateItem} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Mục / Cấp Bậc / Cảnh Giới *
                </label>
                <input
                  type="text"
                  required
                  value={itemName}
                  onChange={(e) => setItemName(e.target.value)}
                  placeholder="Ví dụ: Kim Đan Kỳ, Ma Pháp Sư Sơ Cấp..."
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Thứ Tự Cấp Bậc (Order Index: 1, 2, 3...)
                </label>
                <input
                  type="number"
                  value={orderIndex}
                  onChange={(e) => setOrderIndex(parseInt(e.target.value) || 0)}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả / Đặc Điểm Của Cấp Bậc Này
                </label>
                <textarea
                  rows={2}
                  value={itemDesc}
                  onChange={(e) => setItemDesc(e.target.value)}
                  placeholder="Ghi chú về sức mạnh, tuổi thọ, dấu hiệu..."
                  className="input-field text-xs"
                />
              </div>
              <div className="flex justify-end space-x-2 pt-2 border-t border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsCreateItemOpen(false)}
                  className="px-3 py-1.5 text-xs text-slate-400 hover:text-white"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="btn-primary text-xs"
                >
                  {isSubmitting ? 'Đang thêm...' : 'Thêm Mục'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- Modal: Sửa Cấp Bậc / Mục --- */}
      {isEditItemOpen && editingItem && activeCategory && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => {
                setIsEditItemOpen(false);
                setEditingItem(null);
              }}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4 flex items-center space-x-2">
              <Edit2 className="w-5 h-5 text-brand-400" />
              <span>Chỉnh Sửa Cấp Bậc / Mục</span>
            </h3>
            <form onSubmit={handleUpdateItem} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Mục / Cấp Bậc *
                </label>
                <input
                  type="text"
                  required
                  value={itemName}
                  onChange={(e) => setItemName(e.target.value)}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Thứ Tự Cấp Bậc (Order Index)
                </label>
                <input
                  type="number"
                  value={orderIndex}
                  onChange={(e) => setOrderIndex(parseInt(e.target.value) || 0)}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả / Đặc Điểm
                </label>
                <textarea
                  rows={2}
                  value={itemDesc}
                  onChange={(e) => setItemDesc(e.target.value)}
                  className="input-field text-xs"
                />
              </div>
              <div className="flex justify-end space-x-2 pt-2 border-t border-slate-800">
                <button
                  type="button"
                  onClick={() => {
                    setIsEditItemOpen(false);
                    setEditingItem(null);
                  }}
                  className="px-3 py-1.5 text-xs text-slate-400 hover:text-white"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="btn-primary text-xs"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Lưu Thay Đổi'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- Modal: Lưu Thành Bộ Khung Mẫu (Framework Template) --- */}
      {isSaveTemplateOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-lg w-full p-6 relative">
            <button
              onClick={() => setIsSaveTemplateOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <div className="flex items-center space-x-2 mb-3">
              <div className="p-2 rounded-xl bg-indigo-500/10 border border-indigo-500/20 text-indigo-400">
                <BookmarkPlus className="w-5 h-5" />
              </div>
              <div>
                <h3 className="text-base font-bold text-white">
                  Lưu Bộ Khung Thiết Lập Thế Giới
                </h3>
                <p className="text-xs text-slate-400">
                  Đóng gói toàn bộ các danh mục và cấp bậc tu luyện hiện tại thành mẫu tái sử dụng
                </p>
              </div>
            </div>

            {tmplSuccessMsg ? (
              <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-sm flex items-center space-x-2 my-4">
                <Check className="w-5 h-5 flex-shrink-0" />
                <span>{tmplSuccessMsg}</span>
              </div>
            ) : (
              <form onSubmit={handleSaveFramework} className="space-y-4 mt-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Tên Bộ Khung Mẫu *
                  </label>
                  <input
                    type="text"
                    required
                    value={tmplName}
                    onChange={(e) => setTmplName(e.target.value)}
                    placeholder="VD: Hệ Thống Tu Tiên Phàm Nhân Lưu, Cyberpunk 2099..."
                    className="input-field"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Mô Tả Bộ Khung
                  </label>
                  <textarea
                    rows={3}
                    value={tmplDesc}
                    onChange={(e) => setTmplDesc(e.target.value)}
                    placeholder="Giới thiệu về các quy luật, sự độc đáo của hệ thống cảnh giới này..."
                    className="input-field text-xs"
                  />
                </div>

                <div className="p-3.5 bg-slate-950/60 rounded-xl border border-slate-800 flex items-center justify-between">
                  <div className="flex items-center space-x-2.5">
                    <Share2 className="w-4 h-4 text-purple-400" />
                    <div>
                      <div className="text-xs font-bold text-slate-200">
                        Chia sẻ Thư Viện Cộng Đồng (Public)
                      </div>
                      <div className="text-[11px] text-slate-500">
                        Cho phép các tác giả khác khám phá và sử dụng bộ khung của bạn
                      </div>
                    </div>
                  </div>

                  <input
                    type="checkbox"
                    checked={tmplIsPublic}
                    onChange={(e) => setTmplIsPublic(e.target.checked)}
                    className="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 bg-slate-900 border-slate-700"
                  />
                </div>

                <div className="flex justify-end space-x-2 pt-3 border-t border-slate-800">
                  <button
                    type="button"
                    onClick={() => setIsSaveTemplateOpen(false)}
                    className="btn-secondary text-xs"
                  >
                    Hủy
                  </button>
                  <button
                    type="submit"
                    disabled={isSubmitting}
                    className="btn-primary text-xs flex items-center space-x-1.5"
                  >
                    <BookmarkPlus className="w-4 h-4" />
                    <span>{isSubmitting ? 'Đang lưu...' : 'Lưu Thành Bộ Khung'}</span>
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
