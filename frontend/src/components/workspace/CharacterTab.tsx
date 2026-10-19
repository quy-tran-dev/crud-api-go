import React, { useState } from 'react';
import type { Character, Faction, Location, SystemCategory } from '../../types';
import { api, ENDPOINTS, handleFormSubmit, confirmAndDelete } from '../../services';
import { Users, Plus, Trash2, X, Sparkles, MapPin, Shield } from 'lucide-react';

interface CharacterTabProps {
  storyID: number;
  characters: Character[];
  factions: Faction[];
  locations: Location[];
  categories: SystemCategory[];
  onRefresh: () => void;
  onDetectedNewItems?: (items: any[]) => void;
}

export const CharacterTab: React.FC<CharacterTabProps> = ({
  storyID,
  characters,
  factions,
  locations,
  categories,
  onRefresh,
  onDetectedNewItems,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [name, setName] = useState('');
  const [alias, setAlias] = useState('');
  const [role, setRole] = useState('Nhân vật chính');
  const [gender, setGender] = useState('Nam');
  const [age, setAge] = useState('18');
  const [factionID, setFactionID] = useState<number | null>(null);
  const [factionRank, setFactionRank] = useState('Đệ tử');
  const [birthplaceID, setBirthplaceID] = useState<number | null>(null);
  const [currentLevel, setCurrentLevel] = useState('');
  const [abilities, setAbilities] = useState('');
  const [originStory, setOriginStory] = useState('');
  const [characterArc, setCharacterArc] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Lấy danh sách gợi ý cảnh giới từ category power_level
  const powerLevelCategory = categories.find((c) => c.code === 'power_level');

  const handleCreate = async (e: React.SubmitEvent) => {
    e.preventDefault();
    await handleFormSubmit({
      setSubmitting: setIsSubmitting,
      action: () =>
        api.post(ENDPOINTS.CHARACTERS.BASE(storyID), {
          name,
          alias,
          role,
          gender,
          age,
          faction_id: factionID || null,
          faction_rank: factionRank,
          birthplace_id: birthplaceID || null,
          current_level: currentLevel,
          abilities,
          origin_story: originStory,
          character_arc: characterArc,
        }),
      onSuccess: (response) => {
        // Nếu backend phát hiện từ mới (Smart Recheck)
        if (response.data.detected_new_items && response.data.detected_new_items.length > 0) {
          if (onDetectedNewItems) {
            onDetectedNewItems(response.data.detected_new_items);
          }
        }

        // Reset form
        setName('');
        setAlias('');
        setAbilities('');
        setOriginStory('');
        setCharacterArc('');
        setFactionID(null);
        setBirthplaceID(null);
        setIsModalOpen(false);
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi tạo nhân vật',
      successMessage: `Đã khởi tạo nhân vật "${name}" thành công!`,
    });
  };

  const handleDelete = async (id: number, charName: string) => {
    await confirmAndDelete({
      title: 'Xóa Nhân Vật',
      confirmMessage: `Bạn có chắc muốn xóa nhân vật "${charName}"? Hành động này không thể hoàn tác.`,
      action: () => api.delete(ENDPOINTS.CHARACTERS.DETAIL(storyID, id)),
      onSuccess: () => {
        onRefresh();
      },
      fallbackErrorMessage: 'Lỗi khi xóa nhân vật',
      successMessage: `Đã xóa nhân vật "${charName}" thành công!`,
    });
  };

  return (
    <div className="p-6 overflow-y-auto h-full flex flex-col">
      <div className="flex items-center justify-between pb-4 border-b border-slate-800 mb-6">
        <div>
          <h3 className="text-xl font-bold text-white flex items-center space-x-2">
            <Users className="w-5 h-5 text-amber-400" />
            <span>Hệ Thống Nhân Vật & Tuyến Truyện</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Gia thế, xuất thân, cảnh giới tu vi và mối quan hệ giữa các nhân vật
          </p>
        </div>

        <button
          onClick={() => setIsModalOpen(true)}
          className="inline-flex items-center space-x-1.5 bg-amber-600 hover:bg-amber-500 text-white text-xs font-semibold px-3 py-2 rounded-xl transition-all shadow"
        >
          <Plus className="w-4 h-4" />
          <span>Thêm Nhân Vật Mới</span>
        </button>
      </div>

      {characters.length === 0 ? (
        <div className="text-center py-20 text-slate-500">
          Chưa có nhân vật nào. Hãy sáng tạo nhân vật chính diện hoặc phản diện đầu tiên!
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {characters.map((c) => {
            const faction = factions.find((f) => f.id === c.faction_id);
            const birthplace = locations.find((l) => l.id === c.birthplace_id);

            return (
              <div
                key={c.id}
                className="bg-slate-900 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between hover:border-amber-500/50 transition-all"
              >
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs font-semibold px-2.5 py-0.5 bg-amber-950/80 text-amber-400 border border-amber-800/60 rounded-md">
                      {c.role || 'Nhân vật'}
                    </span>
                    {c.current_level && (
                      <span className="text-xs font-medium text-brand-300 bg-brand-950 px-2 py-0.5 rounded border border-brand-800/50">
                        {c.current_level}
                      </span>
                    )}
                  </div>

                  <h4 className="text-lg font-bold text-white mb-0.5">
                    {c.name}
                    {c.alias && (
                      <span className="text-xs font-normal text-slate-400 ml-2">
                        ({c.alias})
                      </span>
                    )}
                  </h4>

                  <div className="text-xs text-slate-500 mb-3">
                    {c.gender} • {c.age} tuổi
                  </div>

                  {/* Gia thế & Nơi sinh */}
                  <div className="space-y-1.5 text-xs text-slate-300 mb-3 bg-slate-950/60 p-2.5 rounded-xl border border-slate-800/80">
                    {faction ? (
                      <div className="flex items-center space-x-1.5 truncate">
                        <Shield className="w-3.5 h-3.5 text-blue-400 flex-shrink-0" />
                        <span className="truncate">
                          {faction.name} ({c.faction_rank || 'Thành viên'})
                        </span>
                      </div>
                    ) : (
                      <div className="text-slate-500 italic">Chưa gia nhập thế lực nào</div>
                    )}

                    {birthplace && (
                      <div className="flex items-center space-x-1.5 truncate">
                        <MapPin className="w-3.5 h-3.5 text-emerald-400 flex-shrink-0" />
                        <span className="truncate">Quê quán: {birthplace.name}</span>
                      </div>
                    )}
                  </div>

                  {c.origin_story && (
                    <p className="text-xs text-slate-400 line-clamp-2 mb-2">
                      <strong className="text-slate-300">Quá khứ:</strong> {c.origin_story}
                    </p>
                  )}
                </div>

                <div className="pt-3 border-t border-slate-800/80 flex items-center justify-between">
                  <span className="text-[11px] text-slate-500 font-mono">ID: #{c.id}</span>
                  <button
                    onClick={() => handleDelete(c.id, c.name)}
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

      {/* Modal Thêm Nhân Vật */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-xl w-full p-6 relative max-h-[90vh] overflow-y-auto">
            <button
              onClick={() => setIsModalOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
            <h3 className="text-base font-bold text-white mb-4 flex items-center space-x-2">
              <Sparkles className="w-5 h-5 text-amber-400" />
              <span>Sáng Tạo Nhân Vật Mới</span>
            </h3>

            <form onSubmit={handleCreate} className="space-y-3.5">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Tên Nhân Vật *
                  </label>
                  <input
                    type="text"
                    required
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Lâm Phong, Tiêu Viêm..."
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Biệt Danh / Đạo Hiệu
                  </label>
                  <input
                    type="text"
                    value={alias}
                    onChange={(e) => setAlias(e.target.value)}
                    placeholder="Hắc Y Kiếm Khách..."
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-sm text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  />
                </div>
              </div>

              <div className="grid grid-cols-3 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Vai Trò
                  </label>
                  <input
                    type="text"
                    value={role}
                    onChange={(e) => setRole(e.target.value)}
                    placeholder="Chính, Phản diện..."
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Giới Tính
                  </label>
                  <select
                    value={gender}
                    onChange={(e) => setGender(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  >
                    <option value="Nam">Nam</option>
                    <option value="Nữ">Nữ</option>
                    <option value="Khác">Khác / Vô tính</option>
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Tuổi Tác
                  </label>
                  <input
                    type="text"
                    value={age}
                    onChange={(e) => setAge(e.target.value)}
                    placeholder="18"
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  />
                </div>
              </div>

              {/* Cảnh giới tu vi (Có gợi ý hoặc tự do) */}
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Cảnh Giới / Cấp Bậc (Gõ từ mới sẽ tự động Recheck)
                </label>
                <div className="flex flex-wrap gap-1.5 mb-2">
                  {powerLevelCategory?.items?.map((item) => (
                    <button
                      type="button"
                      key={item.id}
                      onClick={() => setCurrentLevel(item.name)}
                      className={`text-[11px] px-2 py-0.5 rounded-md border transition-all ${
                        currentLevel === item.name
                          ? 'bg-brand-600 text-white border-brand-400'
                          : 'bg-slate-800 text-slate-400 border-slate-700 hover:text-white'
                      }`}
                    >
                      {item.name}
                    </button>
                  ))}
                </div>
                <input
                  type="text"
                  value={currentLevel}
                  onChange={(e) => setCurrentLevel(e.target.value)}
                  placeholder="Chọn từ danh sách trên hoặc tự gõ cảnh giới mới (ví dụ: Ngưng Khí Kỳ)..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                />
              </div>

              {/* Gia tộc & Nơi sinh */}
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Gia Tộc / Môn Phái
                  </label>
                  <select
                    value={factionID || ''}
                    onChange={(e) => setFactionID(e.target.value ? parseInt(e.target.value) : null)}
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  >
                    <option value="">Tán tu / Chưa gia nhập</option>
                    {factions.map((f) => (
                      <option key={f.id} value={f.id}>
                        {f.name}
                      </option>
                    ))}
                  </select>
                  {factionID && (
                    <input
                      type="text"
                      value={factionRank}
                      onChange={(e) => setFactionRank(e.target.value)}
                      placeholder="Thân phận: Đệ tử, Trưởng lão..."
                      className="mt-1.5 w-full px-2.5 py-1.5 bg-slate-950 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:ring-1 focus:ring-amber-500"
                    />
                  )}
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                    Nơi Sinh / Quê Quán
                  </label>
                  <select
                    value={birthplaceID || ''}
                    onChange={(e) => setBirthplaceID(e.target.value ? parseInt(e.target.value) : null)}
                    className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  >
                    <option value="">Chưa rõ</option>
                    {locations.map((l) => (
                      <option key={l.id} value={l.id}>
                        {l.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Quá Khứ / Gia Thế (Origin Story)
                </label>
                <textarea
                  rows={2}
                  value={originStory}
                  onChange={(e) => setOriginStory(e.target.value)}
                  placeholder="Gia tộc bị diệt môn, nhặt được chí bảo từ nhỏ..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tuyến Phát Triển Tâm Lý (Character Arc)
                </label>
                <input
                  type="text"
                  value={characterArc}
                  onChange={(e) => setCharacterArc(e.target.value)}
                  placeholder="Từ non nớt ➔ Báo thù ➔ Đại nghĩa thiên hạ..."
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-xs text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
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
                  className="bg-amber-600 hover:bg-amber-500 text-white text-xs font-semibold px-4 py-2 rounded-xl transition-all"
                >
                  {isSubmitting ? 'Đang lưu...' : 'Tạo Nhân Vật'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
