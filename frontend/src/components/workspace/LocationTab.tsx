import React, { useState } from 'react';
import type { Location } from '../../types';
import { api, ENDPOINTS, handleFormSubmit, confirmAndDelete } from '../../services';
import { MapPin, Plus, Trash2, X, CornerDownRight } from 'lucide-react';

interface LocationTabProps {
  storyID: number;
  locations: Location[];
  onRefresh: () => void;
}

export const LocationTab: React.FC<LocationTabProps> = ({
  storyID,
  locations,
  onRefresh,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [name, setName] = useState('');
  const [type, setType] = useState('Đại Lục');
  const [parentID, setParentID] = useState<number | null>(null);
  const [dangerLevel, setDangerLevel] = useState('Thấp');
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleCreate = async (e: React.SubmitEvent) => {
    e.preventDefault();
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.LOCATIONS.BASE(storyID), {
          name,
          type,
          parent_id: parentID || null,
          danger_level: dangerLevel,
          description,
        }),
      onSuccess: () => {
        setName('');
        setDescription('');
        setParentID(null);
        setIsModalOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi tạo địa danh',
      successMessage: `Đã tạo địa danh "${name}" thành công!`,
    });
  };

  const handleDelete = async (id: number, locName: string) => {
    await confirmAndDelete({
      title: 'Xóa Địa Danh',
      confirmMessage: `Bạn có chắc muốn xóa địa danh "${locName}"? Các địa danh con trực thuộc cũng có thể bị ảnh hưởng.`,
      action: () => api.delete(ENDPOINTS.LOCATIONS.DETAIL(storyID, id)),
      onSuccess: () => {
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi xóa địa danh',
      successMessage: `Đã xóa địa danh "${locName}" thành công!`,
    });
  };

  return (
    <div className="p-6 overflow-y-auto h-full flex flex-col">
      <div className="flex items-center justify-between pb-4 border-b border-slate-800 mb-6">
        <div>
          <h3 className="text-xl font-bold text-white flex items-center space-x-2">
            <MapPin className="w-5 h-5 text-emerald-400" />
            <span>Bản Đồ & Địa Danh (Cây Phân Cấp)</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Phân cấp địa lý từ Đại lục, Quốc gia đến Tông môn và Bí cảnh
          </p>
        </div>

        <button
          onClick={() => setIsModalOpen(true)}
          className="inline-flex items-center space-x-1.5 bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold px-3 py-2 rounded-xl transition-all shadow"
        >
          <Plus className="w-4 h-4" />
          <span>Thêm Địa Danh Mới</span>
        </button>
      </div>

      {/* Danh sách địa danh */}
      {locations.length === 0 ? (
        <div className="text-center py-20 text-slate-500">
          Chưa có địa danh nào. Hãy tạo Đại lục hoặc Sơn môn đầu tiên!
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {locations.map((loc) => (
            <div
              key={loc.id}
              className="bg-slate-900 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between hover:border-emerald-500/50 transition-all"
            >
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-semibold px-2.5 py-0.5 bg-emerald-950/80 text-emerald-400 border border-emerald-800/60 rounded-md">
                    {loc.type || 'Địa danh'}
                  </span>
                  {loc.danger_level && (
                    <span className="text-[11px] text-slate-400 bg-slate-800 px-2 py-0.5 rounded border border-slate-700">
                      Nguy hiểm: {loc.danger_level}
                    </span>
                  )}
                </div>

                <h4 className="text-base font-bold text-white mb-1">{loc.name}</h4>
                <p className="text-xs text-slate-400 line-clamp-2 mb-3">
                  {loc.description || 'Chưa có mô tả địa hình...'}
                </p>

                {loc.parent_id && (
                  <div className="flex items-center space-x-1 text-[11px] text-emerald-400/80 font-mono mb-2">
                    <CornerDownRight className="w-3.5 h-3.5" />
                    <span>Thuộc Map cha ID: #{loc.parent_id}</span>
                  </div>
                )}
              </div>

              <div className="pt-3 border-t border-slate-800/80 flex items-center justify-between">
                <span className="text-[11px] text-slate-500 font-mono">ID: #{loc.id}</span>
                <button
                  onClick={() => handleDelete(loc.id, loc.name)}
                  className="p-1.5 text-slate-500 hover:text-red-400 hover:bg-slate-800 rounded-lg transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal Thêm Địa Danh */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => setIsModalOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4">Thêm Địa Danh Mới</h3>
            <form onSubmit={handleCreate} className="space-y-3.5">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Địa Danh *
                </label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Ví dụ: Càn Khôn Đại Lục, Vân Lam Sơn Môn, Vạn Thú Cấm Địa..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Phân Loại
                  </label>
                  <input
                    type="text"
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                    placeholder="Đại Lục, Bí Cảnh..."
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Mức Độ Nguy Hiểm
                  </label>
                  <select
                    value={dangerLevel}
                    onChange={(e) => setDangerLevel(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="Thấp">Thấp</option>
                    <option value="Trung bình">Trung bình</option>
                    <option value="Cao">Cao</option>
                    <option value="Tử địa">Tử địa / Cực nguy</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Thuộc Về Địa Danh Cha (Để tạo cây phân cấp)
                </label>
                <select
                  value={parentID || ''}
                  onChange={(e) => setParentID(e.target.value ? parseInt(e.target.value) : null)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                >
                  <option value="">Không có (Là Map gốc cao nhất)</option>
                  {locations.map((l) => (
                    <option key={l.id} value={l.id}>
                      {l.name} ({l.type})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả Địa Hình & Tài Nguyên
                </label>
                <textarea
                  rows={2}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Khí hậu, địa mạo, linh khí nồng đậm..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
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
                  className="bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold px-4 py-2 rounded-xl transition-all"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Thêm Địa Danh'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
