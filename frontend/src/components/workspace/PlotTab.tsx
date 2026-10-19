import React, { useState } from 'react';
import type { Plot, Location, Character, Faction, Item } from '../../types';
import { api, ENDPOINTS, handleFormSubmit, confirmAndDelete } from '../../services';
import { Film, Plus, Trash2, X, Clock, MapPin, Users, Shield, Sparkles } from 'lucide-react';

interface PlotTabProps {
  storyID: number;
  plots: Plot[];
  locations: Location[];
  characters: Character[];
  factions: Faction[];
  items: Item[];
  onRefresh: () => void;
}

export const PlotTab: React.FC<PlotTabProps> = ({
  storyID,
  plots,
  locations,
  characters,
  factions,
  items,
  onRefresh,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [type, setType] = useState('main_plot');
  const [orderIndex, setOrderIndex] = useState(plots.length + 1);
  const [locationID, setLocationID] = useState<number | null>(null);
  const [summary, setSummary] = useState('');
  const [consequences, setConsequences] = useState('');
  const [selectedCharacters, setSelectedCharacters] = useState<number[]>([]);
  const [selectedFactions, setSelectedFactions] = useState<number[]>([]);
  const [selectedItems, setSelectedItems] = useState<number[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Sắp xếp tình tiết theo order_index tăng dần
  const sortedPlots = [...plots].sort((a, b) => a.order_index - b.order_index);

  const handleCreate = async (e: React.SubmitEvent) => {
    e.preventDefault();
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.PLOTS.BASE(storyID), {
          title,
          type,
          order_index: Number(orderIndex),
          location_id: locationID || null,
          summary,
          consequences,
          participating_character_ids: selectedCharacters,
          involved_faction_ids: selectedFactions,
          featured_item_ids: selectedItems,
        }),
      onSuccess: () => {
        // Reset form
        setTitle('');
        setType('main_plot');
        setOrderIndex(plots.length + 2);
        setLocationID(null);
        setSummary('');
        setConsequences('');
        setSelectedCharacters([]);
        setSelectedFactions([]);
        setSelectedItems([]);
        setIsModalOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi tạo tình tiết',
      successMessage: `Đã tạo tình tiết "${title}" thành công!`,
    });
  };

  const handleDelete = async (id: number, plotTitle: string) => {
    await confirmAndDelete({
      title: 'Xóa Tình Tiết',
      confirmMessage: `Bạn có chắc muốn xóa tình tiết "${plotTitle}"?`,
      action: () => api.delete(ENDPOINTS.PLOTS.DETAIL(storyID, id)),
      onSuccess: () => {
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi xóa tình tiết',
      successMessage: `Đã xóa tình tiết "${plotTitle}" thành công!`,
    });
  };

  const getTypeBadge = (plotType: string) => {
    switch (plotType) {
      case 'main_plot':
        return <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-rose-500/20 text-rose-300 border border-rose-500/30">Mạch Chính</span>;
      case 'sub_plot':
        return <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-purple-500/20 text-purple-300 border border-purple-500/30">Nhánh Phụ</span>;
      case 'side_story':
        return <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-cyan-500/20 text-cyan-300 border border-cyan-500/30">Ngoại Truyện</span>;
      default:
        return <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-slate-700 text-slate-300">{plotType}</span>;
    }
  };

  const toggleSelect = (id: number, list: number[], setList: React.Dispatch<React.SetStateAction<number[]>>) => {
    if (list.includes(id)) {
      setList(list.filter((x) => x !== id));
    } else {
      setList([...list, id]);
    }
  };

  return (
    <div className="p-6 overflow-y-auto h-full flex flex-col">
      {/* Header */}
      <div className="flex items-center justify-between pb-5 border-b border-slate-800">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center space-x-2">
            <Film className="w-5 h-5 text-purple-400" />
            <span>Dòng Thời Gian & Tình Tiết (Timeline Plots)</span>
          </h2>
          <p className="text-xs text-slate-400 mt-1">
            Quản lý các sự kiện, nút thắt kịch bản theo thứ tự thời gian và liên kết nhân vật, thế lực, vật phẩm
          </p>
        </div>
        <button
          onClick={() => {
            setOrderIndex(plots.length + 1);
            setIsModalOpen(true);
          }}
          className="btn-primary flex items-center space-x-1.5"
        >
          <Plus className="w-4 h-4" />
          <span>Thêm Tình Tiết</span>
        </button>
      </div>

      {/* Content Timeline */}
      <div className="mt-6 flex-1">
        {sortedPlots.length === 0 ? (
          <div className="bg-slate-900 border border-dashed border-slate-800 rounded-2xl p-12 text-center text-slate-400">
            <Clock className="w-12 h-12 mx-auto mb-3 text-slate-600" />
            <p className="font-semibold text-slate-300">Chưa có tình tiết hoặc biến cố nào</p>
            <p className="text-xs mt-1 text-slate-500">
              Hãy thêm sự kiện mở đầu cho câu chuyện của bạn
            </p>
          </div>
        ) : (
          <div className="relative border-l-2 border-purple-500/30 ml-4 pl-6 space-y-6">
            {sortedPlots.map((plot) => {
              const loc = locations.find((l) => l.id === plot.location_id);
              const partChars = characters.filter((c) =>
                plot.participating_character_ids?.includes(c.id)
              );
              const invFactions = factions.filter((f) =>
                plot.involved_faction_ids?.includes(f.id)
              );
              const featItems = items.filter((i) =>
                plot.featured_item_ids?.includes(i.id)
              );

              return (
                <div key={plot.id} className="relative group">
                  {/* Timeline node */}
                  <div className="absolute -left-[31px] top-1.5 w-4 h-4 rounded-full bg-purple-500 border-4 border-slate-950 shadow-sm" />

                  {/* Card Content */}
                  <div className="bg-slate-900 border border-slate-800 rounded-xl p-4 hover:border-slate-700 transition">
                    <div className="flex items-start justify-between">
                      <div className="space-y-1">
                        <div className="flex items-center space-x-2">
                          <span className="text-xs font-mono font-bold text-purple-400 bg-purple-500/10 px-2 py-0.5 rounded">
                            #{plot.order_index}
                          </span>
                          <h3 className="font-bold text-white text-base">{plot.title}</h3>
                          {getTypeBadge(plot.type)}
                        </div>

                        {loc && (
                          <div className="flex items-center text-xs text-emerald-400 space-x-1 pt-0.5">
                            <MapPin className="w-3.5 h-3.5" />
                            <span>{loc.name}</span>
                          </div>
                        )}
                      </div>

                      <button
                        onClick={() => handleDelete(plot.id, plot.title)}
                        className="text-slate-500 hover:text-rose-400 p-1 rounded-lg hover:bg-slate-800 transition"
                        title="Xóa tình tiết"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>

                    {plot.summary && (
                      <p className="text-sm text-slate-300 mt-3 whitespace-pre-line leading-relaxed">
                        {plot.summary}
                      </p>
                    )}

                    {plot.consequences && (
                      <div className="mt-3 p-2.5 rounded-lg bg-amber-500/10 border border-amber-500/20 text-xs text-amber-300 leading-normal">
                        <span className="font-semibold text-amber-200">Hậu quả / Biến cố tiếp theo: </span>
                        {plot.consequences}
                      </div>
                    )}

                    {/* Metadata tags */}
                    <div className="mt-4 pt-3 border-t border-slate-800/80 flex flex-wrap gap-2">
                      {partChars.length > 0 && (
                        <div className="flex items-center space-x-1 text-[11px] text-amber-400 bg-amber-500/10 px-2 py-1 rounded-md">
                          <Users className="w-3 h-3" />
                          <span>
                            Nhân vật: {partChars.map((c) => c.name).join(', ')}
                          </span>
                        </div>
                      )}

                      {invFactions.length > 0 && (
                        <div className="flex items-center space-x-1 text-[11px] text-blue-400 bg-blue-500/10 px-2 py-1 rounded-md">
                          <Shield className="w-3 h-3" />
                          <span>
                            Thế lực: {invFactions.map((f) => f.name).join(', ')}
                          </span>
                        </div>
                      )}

                      {featItems.length > 0 && (
                        <div className="flex items-center space-x-1 text-[11px] text-rose-400 bg-rose-500/10 px-2 py-1 rounded-md">
                          <Sparkles className="w-3 h-3" />
                          <span>
                            Vật phẩm: {featItems.map((i) => i.name).join(', ')}
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Modal Thêm Tình Tiết */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-2xl p-6 shadow-2xl max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between pb-4 border-b border-slate-800">
              <h3 className="text-lg font-bold text-white flex items-center space-x-2">
                <Film className="w-5 h-5 text-purple-400" />
                <span>Thêm Tình Tiết / Sự Kiện Mới</span>
              </h3>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreate} className="mt-4 space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="md:col-span-2">
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Tiêu đề sự kiện *
                  </label>
                  <input
                    type="text"
                    required
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="VD: Trận chiến Vạn Hồn Nhai, Đại hội Võ lâm..."
                    className="input-field"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Thứ tự diễn biến
                  </label>
                  <input
                    type="number"
                    value={orderIndex}
                    onChange={(e) => setOrderIndex(Number(e.target.value))}
                    className="input-field"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Phân loại
                  </label>
                  <select
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                    className="input-field"
                  >
                    <option value="main_plot">Mạch Chính (Main Plot)</option>
                    <option value="sub_plot">Nhánh Phụ (Sub Plot)</option>
                    <option value="side_story">Ngoại Truyện (Side Story)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Địa điểm xảy ra
                  </label>
                  <select
                    value={locationID || ''}
                    onChange={(e) =>
                      setLocationID(e.target.value ? Number(e.target.value) : null)
                    }
                    className="input-field"
                  >
                    <option value="">-- Không xác định --</option>
                    {locations.map((loc) => (
                      <option key={loc.id} value={loc.id}>
                        {loc.name} ({loc.type})
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 mb-1">
                  Diễn biến chi tiết
                </label>
                <textarea
                  rows={3}
                  value={summary}
                  onChange={(e) => setSummary(e.target.value)}
                  placeholder="Mô tả sự việc khởi đầu, cao trào và kết thúc..."
                  className="input-field"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 mb-1">
                  Hậu quả / Biến cố tiếp theo
                </label>
                <textarea
                  rows={2}
                  value={consequences}
                  onChange={(e) => setConsequences(e.target.value)}
                  placeholder="Những thay đổi sau tình tiết (ví dụ: nhân vật bị trục xuất, bang hội sụp đổ...)"
                  className="input-field"
                />
              </div>

              {/* Nhân vật tham gia */}
              {characters.length > 0 && (
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1.5">
                    Nhân vật tham gia
                  </label>
                  <div className="flex flex-wrap gap-2 max-h-28 overflow-y-auto p-2 bg-slate-950/60 rounded-xl border border-slate-800">
                    {characters.map((c) => {
                      const isSelected = selectedCharacters.includes(c.id);
                      return (
                        <button
                          key={c.id}
                          type="button"
                          onClick={() => toggleSelect(c.id, selectedCharacters, setSelectedCharacters)}
                          className={`text-xs px-2.5 py-1 rounded-lg border transition ${
                            isSelected
                              ? 'bg-amber-500/20 border-amber-500/40 text-amber-300 font-semibold'
                              : 'bg-slate-800 border-slate-700 text-slate-400 hover:text-slate-200'
                          }`}
                        >
                          {c.name}
                        </button>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Thế lực dính líu */}
              {factions.length > 0 && (
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1.5">
                    Thế lực liên quan
                  </label>
                  <div className="flex flex-wrap gap-2 max-h-28 overflow-y-auto p-2 bg-slate-950/60 rounded-xl border border-slate-800">
                    {factions.map((f) => {
                      const isSelected = selectedFactions.includes(f.id);
                      return (
                        <button
                          key={f.id}
                          type="button"
                          onClick={() => toggleSelect(f.id, selectedFactions, setSelectedFactions)}
                          className={`text-xs px-2.5 py-1 rounded-lg border transition ${
                            isSelected
                              ? 'bg-blue-500/20 border-blue-500/40 text-blue-300 font-semibold'
                              : 'bg-slate-800 border-slate-700 text-slate-400 hover:text-slate-200'
                          }`}
                        >
                          {f.name}
                        </button>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Vật phẩm xuất hiện */}
              {items.length > 0 && (
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1.5">
                    Vật phẩm xuất hiện
                  </label>
                  <div className="flex flex-wrap gap-2 max-h-28 overflow-y-auto p-2 bg-slate-950/60 rounded-xl border border-slate-800">
                    {items.map((i) => {
                      const isSelected = selectedItems.includes(i.id);
                      return (
                        <button
                          key={i.id}
                          type="button"
                          onClick={() => toggleSelect(i.id, selectedItems, setSelectedItems)}
                          className={`text-xs px-2.5 py-1 rounded-lg border transition ${
                            isSelected
                              ? 'bg-rose-500/20 border-rose-500/40 text-rose-300 font-semibold'
                              : 'bg-slate-800 border-slate-700 text-slate-400 hover:text-slate-200'
                          }`}
                        >
                          {i.name}
                        </button>
                      );
                    })}
                  </div>
                </div>
              )}

              <div className="pt-4 border-t border-slate-800 flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 rounded-xl text-xs font-semibold text-slate-400 hover:bg-slate-800"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="btn-primary"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Lưu Tình Tiết'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
