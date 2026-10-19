import React, { useState } from 'react';
import type { Faction, Location } from '../../types';
import { api, ENDPOINTS, handleFormSubmit, confirmAndDelete } from '../../services';
import { Shield, Plus, Trash2, X, MapPin } from 'lucide-react';

interface FactionTabProps {
  storyID: number;
  factions: Faction[];
  locations: Location[];
  onRefresh: () => void;
}

export const FactionTab: React.FC<FactionTabProps> = ({
  storyID,
  factions,
  locations,
  onRefresh,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [name, setName] = useState('');
  const [type, setType] = useState('Tông Môn Chính Đạo');
  const [powerTier, setPowerTier] = useState('Nhất Lưu');
  const [headquarterID, setHeadquarterID] = useState<number | null>(null);
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleCreate = async (e: React.SubmitEvent) => {
    e.preventDefault();
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.FACTIONS.BASE(storyID), {
          name,
          type,
          power_tier: powerTier,
          headquarter_id: headquarterID || null,
          description,
        }),
      onSuccess: () => {
        setName('');
        setDescription('');
        setHeadquarterID(null);
        setIsModalOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi tạo thế lực',
      successMessage: `Đã tạo thế lực "${name}" thành công!`,
    });
  };

  const handleDelete = async (id: number, facName: string) => {
    await confirmAndDelete({
      title: 'Xóa Thế Lực',
      confirmMessage: `Bạn có chắc muốn xóa thế lực "${facName}"?`,
      action: () => api.delete(ENDPOINTS.FACTIONS.DETAIL(storyID, id)),
      onSuccess: () => {
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi xóa thế lực',
      successMessage: `Đã xóa thế lực "${facName}" thành công!`,
    });
  };

  return (
    <div className="p-6 overflow-y-auto h-full flex flex-col">
      <div className="flex items-center justify-between pb-4 border-b border-slate-800 mb-6">
        <div>
          <h3 className="text-xl font-bold text-white flex items-center space-x-2">
            <Shield className="w-5 h-5 text-blue-400" />
            <span>Thế Lực, Tông Môn & Gia Tộc</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Các phe phái quyền lực, gia tộc thế gia và tổ chức ngầm
          </p>
        </div>

        <button
          onClick={() => setIsModalOpen(true)}
          className="inline-flex items-center space-x-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold px-3 py-2 rounded-xl transition-all shadow"
        >
          <Plus className="w-4 h-4" />
          <span>Thêm Thế Lực Mới</span>
        </button>
      </div>

      {factions.length === 0 ? (
        <div className="text-center py-20 text-slate-500">
          Chưa có thế lực nào. Hãy khai sáng môn phái hoặc gia tộc đầu tiên!
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {factions.map((f) => {
            const hqLocation = locations.find((l) => l.id === f.headquarter_id);

            return (
              <div
                key={f.id}
                className="bg-slate-900 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between hover:border-blue-500/50 transition-all"
              >
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs font-semibold px-2.5 py-0.5 bg-blue-950/80 text-blue-400 border border-blue-800/60 rounded-md">
                      {f.type || 'Thế lực'}
                    </span>
                    {f.power_tier && (
                      <span className="text-[11px] text-slate-300 bg-slate-800 px-2 py-0.5 rounded border border-slate-700">
                        {f.power_tier}
                      </span>
                    )}
                  </div>

                  <h4 className="text-base font-bold text-white mb-1">{f.name}</h4>
                  <p className="text-xs text-slate-400 line-clamp-3 mb-3">
                    {f.description || 'Chưa có mô tả bối cảnh thế lực...'}
                  </p>

                  {hqLocation && (
                    <div className="flex items-center space-x-1.5 text-xs text-slate-400 bg-slate-950/70 p-2 rounded-lg border border-slate-800 mb-2">
                      <MapPin className="w-3.5 h-3.5 text-emerald-400 flex-shrink-0" />
                      <span className="truncate">Trụ sở: {hqLocation.name}</span>
                    </div>
                  )}
                </div>

                <div className="pt-3 border-t border-slate-800/80 flex items-center justify-between">
                  <span className="text-[11px] text-slate-500 font-mono">ID: #{f.id}</span>
                  <button
                    onClick={() => handleDelete(f.id, f.name)}
                    className="p-1.5 text-slate-500 hover:text-red-400 hover:bg-slate-800 rounded-lg transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Modal Thêm Thế Lực */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative">
            <button
              onClick={() => setIsModalOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4">Thêm Thế Lực Mới</h3>
            <form onSubmit={handleCreate} className="space-y-3.5">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Thế Lực / Môn Phái / Gia Tộc *
                </label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Ví dụ: Vân Lam Tông, Tiêu Gia, Hồn Điện..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
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
                    placeholder="Tông Môn, Gia Tộc..."
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Cấp Bậc / Quy Mô
                  </label>
                  <input
                    type="text"
                    value={powerTier}
                    onChange={(e) => setPowerTier(e.target.value)}
                    placeholder="Nhất Lưu, Ẩn Thế..."
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Trụ Sở Đóng Quân Tại Địa Danh Nào
                </label>
                <select
                  value={headquarterID || ''}
                  onChange={(e) => setHeadquarterID(e.target.value ? parseInt(e.target.value) : null)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Chưa xác định (Bản nháp tự do)</option>
                  {locations.map((l) => (
                    <option key={l.id} value={l.id}>
                      {l.name} ({l.type})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Mô Tả & Lịch Sử Tông Phái
                </label>
                <textarea
                  rows={2}
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Người sáng lập, công pháp trấn phái, tôn chỉ hoạt động..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
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
                  className="bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold px-4 py-2 rounded-xl transition-all"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Thêm Thế Lực'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
