import React, { useState } from 'react';
import type { Item, Character, Faction, Location } from '../../types';
import { api, ENDPOINTS, handleFormSubmit, confirmAndDelete } from '../../services';
import { Sword, Plus, Trash2, X, Search, Sparkles, User, Shield, MapPin, Tag } from 'lucide-react';

interface ItemTabProps {
  storyID: number;
  items: Item[];
  characters: Character[];
  factions: Faction[];
  locations: Location[];
  categories: Record<string, string[]>;
  onRefresh: () => void;
}

export const ItemTab: React.FC<ItemTabProps> = ({
  storyID,
  items,
  characters,
  factions,
  locations,
  categories,
  onRefresh,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [typeFilter, setTypeFilter] = useState('ALL');

  // Form states
  const [name, setName] = useState('');
  const [type, setType] = useState('Vũ khí');
  const [grade, setGrade] = useState('');
  const [funcDesc, setFuncDesc] = useState('');
  const [usageCondition, setUsageCondition] = useState('');
  const [originLocationID, setOriginLocationID] = useState<number | null>(null);
  const [ownerCharacterID, setOwnerCharacterID] = useState<number | null>(null);
  const [ownerFactionID, setOwnerFactionID] = useState<number | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Lấy danh sách item types từ categories hệ thống hoặc mặc định
  const itemTypeOptions = categories['Loại Trang Bị'] || [
    'Vũ khí',
    'Pháp bảo',
    'Đan dược',
    'Manh mối',
    'Công pháp',
    'Linh thảo',
  ];

  // Gợi ý phẩm giai phổ biến
  const gradeSuggestions = [
    'Phàm giai',
    'Hoàng giai',
    'Huyền giai',
    'Địa giai',
    'Thiên giai',
    'Thần cấp',
    'Thượng phẩm',
    'Hạ phẩm',
  ];

  const handleCreate = async (e: React.SubmitEvent) => {
    e.preventDefault();
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.ITEMS.BASE(storyID), {
          name,
          type,
          grade,
          function: funcDesc,
          usage_condition: usageCondition,
          origin_location_id: originLocationID || null,
          owner_character_id: ownerCharacterID || null,
          owner_faction_id: ownerFactionID || null,
        }),
      onSuccess: () => {
        // Reset
        setName('');
        setType('Vũ khí');
        setGrade('');
        setFuncDesc('');
        setUsageCondition('');
        setOriginLocationID(null);
        setOwnerCharacterID(null);
        setOwnerFactionID(null);
        setIsModalOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi tạo vật phẩm',
      successMessage: `Đã tạo vật phẩm "${name}" thành công!`,
    });
  };

  const handleDelete = async (id: number, itemName: string) => {
    await confirmAndDelete({
      title: 'Xóa Vật Phẩm',
      confirmMessage: `Bạn có chắc muốn xóa vật phẩm "${itemName}"?`,
      action: () => api.delete(ENDPOINTS.ITEMS.DETAIL(storyID, id)),
      onSuccess: () => {
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi xóa vật phẩm',
      successMessage: `Đã xóa vật phẩm "${itemName}" thành công!`,
    });
  };

  // Filter items
  const filteredItems = items.filter((item) => {
    const matchesSearch =
      item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (item.grade && item.grade.toLowerCase().includes(searchTerm.toLowerCase()));
    const matchesType = typeFilter === 'ALL' || item.type === typeFilter;
    return matchesSearch && matchesType;
  });

  return (
    <div className="p-6 overflow-y-auto h-full flex flex-col">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-5 border-b border-slate-800 gap-4">
        <div>
          <h2 className="text-xl font-bold text-white flex items-center space-x-2">
            <Sword className="w-5 h-5 text-rose-400" />
            <span>Vật Phẩm, Pháp Bảo & Phẩm Cấp</span>
          </h2>
          <p className="text-xs text-slate-400 mt-1">
            Quản lý vũ khí, pháp bảo, đan dược, manh mối cùng hệ thống phân cấp phẩm giai
          </p>
        </div>

        <button
          onClick={() => setIsModalOpen(true)}
          className="btn-primary flex items-center space-x-1.5 self-start sm:self-auto"
        >
          <Plus className="w-4 h-4" />
          <span>Thêm Vật Phẩm</span>
        </button>
      </div>

      {/* Filters */}
      <div className="mt-4 flex flex-col sm:flex-row items-center gap-3">
        <div className="relative flex-1 w-full">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-3" />
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Tìm kiếm vật phẩm, phẩm giai..."
            className="input-field pl-9 text-xs"
          />
        </div>

        <div className="flex items-center space-x-2 w-full sm:w-auto">
          <Tag className="w-4 h-4 text-slate-400" />
          <select
            value={typeFilter}
            onChange={(e) => setTypeFilter(e.target.value)}
            className="input-field text-xs py-2"
          >
            <option value="ALL">Tất cả chủng loại</option>
            {itemTypeOptions.map((opt) => (
              <option key={opt} value={opt}>
                {opt}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Grid Cards */}
      <div className="mt-6 flex-1">
        {filteredItems.length === 0 ? (
          <div className="bg-slate-900 border border-dashed border-slate-800 rounded-2xl p-12 text-center text-slate-400">
            <Sword className="w-12 h-12 mx-auto mb-3 text-slate-600" />
            <p className="font-semibold text-slate-300">Không tìm thấy vật phẩm nào</p>
            <p className="text-xs mt-1 text-slate-500">
              Hãy thêm thần binh lợi khí hoặc đan dược cho nhân vật
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredItems.map((item) => {
              const char = characters.find((c) => c.id === item.owner_character_id);
              const fac = factions.find((f) => f.id === item.owner_faction_id);
              const loc = locations.find((l) => l.id === item.origin_location_id);

              return (
                <div
                  key={item.id}
                  className="bg-slate-900 border border-slate-800 rounded-2xl p-4 hover:border-slate-700 transition flex flex-col justify-between"
                >
                  <div>
                    {/* Header card */}
                    <div className="flex items-start justify-between">
                      <div>
                        <div className="flex items-center space-x-2">
                          <h3 className="font-bold text-white text-base">{item.name}</h3>
                        </div>
                        <div className="flex items-center space-x-2 mt-1.5 flex-wrap gap-1">
                          <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-slate-800 text-slate-300 border border-slate-700">
                            {item.type || 'Chưa phân loại'}
                          </span>
                          {item.grade && (
                            <span className="px-2 py-0.5 rounded text-[11px] font-bold bg-amber-500/20 text-amber-300 border border-amber-500/30 flex items-center space-x-1">
                              <Sparkles className="w-3 h-3 text-amber-400" />
                              <span>{item.grade}</span>
                            </span>
                          )}
                        </div>
                      </div>

                      <button
                        onClick={() => handleDelete(item.id, item.name)}
                        className="text-slate-500 hover:text-rose-400 p-1 rounded-lg hover:bg-slate-800 transition"
                        title="Xóa vật phẩm"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>

                    {/* Function */}
                    {item.function && (
                      <div className="mt-3 text-xs text-slate-300 leading-relaxed bg-slate-950/40 p-2.5 rounded-xl border border-slate-800/60">
                        <span className="font-semibold text-slate-400">Công năng: </span>
                        {item.function}
                      </div>
                    )}

                    {/* Usage condition */}
                    {item.usage_condition && (
                      <div className="mt-2 text-xs text-rose-300/90 leading-relaxed bg-rose-500/5 p-2.5 rounded-xl border border-rose-500/20">
                        <span className="font-semibold text-rose-300">Điều kiện dùng: </span>
                        {item.usage_condition}
                      </div>
                    )}
                  </div>

                  {/* Owner & Origin info */}
                  <div className="mt-4 pt-3 border-t border-slate-800 text-[11px] space-y-1 text-slate-400">
                    {char && (
                      <div className="flex items-center space-x-1.5 text-amber-400">
                        <User className="w-3.5 h-3.5" />
                        <span>Sở hữu: <strong className="font-medium text-slate-200">{char.name}</strong></span>
                      </div>
                    )}
                    {fac && (
                      <div className="flex items-center space-x-1.5 text-blue-400">
                        <Shield className="w-3.5 h-3.5" />
                        <span>Trấn bảo môn phái: <strong className="font-medium text-slate-200">{fac.name}</strong></span>
                      </div>
                    )}
                    {loc && (
                      <div className="flex items-center space-x-1.5 text-emerald-400">
                        <MapPin className="w-3.5 h-3.5" />
                        <span>Xuất xứ: <strong className="font-medium text-slate-200">{loc.name}</strong></span>
                      </div>
                    )}
                    {!char && !fac && !loc && (
                      <span className="italic text-slate-600">Chưa gán chủ sở hữu</span>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Modal Thêm Vật Phẩm */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-xl p-6 shadow-2xl max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between pb-4 border-b border-slate-800">
              <h3 className="text-lg font-bold text-white flex items-center space-x-2">
                <Sword className="w-5 h-5 text-rose-400" />
                <span>Thêm Vật Phẩm & Phẩm Cấp Mới</span>
              </h3>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-slate-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreate} className="mt-4 space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Tên vật phẩm *
                  </label>
                  <input
                    type="text"
                    required
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="VD: Trảm Long Kiếm, Tụ Khí Đan..."
                    className="input-field"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Loại vật phẩm
                  </label>
                  <select
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                    className="input-field"
                  >
                    {itemTypeOptions.map((opt) => (
                      <option key={opt} value={opt}>
                        {opt}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              {/* Phẩm cấp / Phẩm giai */}
              <div>
                <div className="flex items-center justify-between mb-1">
                  <label className="text-xs font-semibold text-slate-300">
                    Phẩm cấp / Phẩm giai
                  </label>
                  <span className="text-[11px] text-slate-500">VD: Địa giai thượng phẩm, Thần cấp...</span>
                </div>
                <input
                  type="text"
                  value={grade}
                  onChange={(e) => setGrade(e.target.value)}
                  placeholder="Nhập phẩm giai hoặc chọn nhanh bên dưới..."
                  className="input-field"
                />
                <div className="flex flex-wrap gap-1.5 mt-2">
                  {gradeSuggestions.map((sug) => (
                    <button
                      key={sug}
                      type="button"
                      onClick={() => setGrade(sug)}
                      className="text-[11px] px-2 py-0.5 rounded bg-slate-800 text-slate-400 hover:text-amber-300 hover:bg-slate-700 transition"
                    >
                      {sug}
                    </button>
                  ))}
                </div>
              </div>

              {/* Sở hữu & Xuất xứ */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Nhân vật sở hữu
                  </label>
                  <select
                    value={ownerCharacterID || ''}
                    onChange={(e) =>
                      setOwnerCharacterID(e.target.value ? Number(e.target.value) : null)
                    }
                    className="input-field text-xs"
                  >
                    <option value="">-- Không --</option>
                    {characters.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Gia tộc / Phái sở hữu
                  </label>
                  <select
                    value={ownerFactionID || ''}
                    onChange={(e) =>
                      setOwnerFactionID(e.target.value ? Number(e.target.value) : null)
                    }
                    className="input-field text-xs"
                  >
                    <option value="">-- Không --</option>
                    {factions.map((f) => (
                      <option key={f.id} value={f.id}>
                        {f.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 mb-1">
                    Nơi xuất xứ / Tìm thấy
                  </label>
                  <select
                    value={originLocationID || ''}
                    onChange={(e) =>
                      setOriginLocationID(e.target.value ? Number(e.target.value) : null)
                    }
                    className="input-field text-xs"
                  >
                    <option value="">-- Không --</option>
                    {locations.map((l) => (
                      <option key={l.id} value={l.id}>
                        {l.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 mb-1">
                  Công năng & Tác dụng
                </label>
                <textarea
                  rows={2}
                  value={funcDesc}
                  onChange={(e) => setFuncDesc(e.target.value)}
                  placeholder="Khả năng đặc biệt, sát thương, trị thương..."
                  className="input-field"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 mb-1">
                  Điều kiện sử dụng / Tác dụng phụ
                </label>
                <textarea
                  rows={2}
                  value={usageCondition}
                  onChange={(e) => setUsageCondition(e.target.value)}
                  placeholder="Yêu cầu tu vi, tác dụng phụ sau khi kích hoạt..."
                  className="input-field"
                />
              </div>

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
                  {isSubmitting ? 'Đang lưu...' : 'Lưu Vật Phẩm'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
