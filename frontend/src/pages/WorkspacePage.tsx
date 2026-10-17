import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import type {
  Story,
  SystemCategory,
  Location,
  Faction,
  Character,
  Plot,
  Item,
} from '../types';
import {
  WorkspaceSidebar,
  type WorkspaceTab,
} from '../components/workspace/WorkspaceSidebar';
import { CategoryTab } from '../components/workspace/CategoryTab';
import { LocationTab } from '../components/workspace/LocationTab';
import { FactionTab } from '../components/workspace/FactionTab';
import { CharacterTab } from '../components/workspace/CharacterTab';
import { PlotTab } from '../components/workspace/PlotTab';
import { ItemTab } from '../components/workspace/ItemTab';
import { SmartRecheckModal } from '../components/workspace/SmartRecheckModal';
import { PublishModal } from '../components/workspace/PublishModal';
import {
  ArrowLeft,
  Sparkles,
  BookOpen,
  CheckCircle2,
  AlertCircle,
  Loader2,
  RefreshCw,
  ShieldCheck,
} from 'lucide-react';

export const WorkspacePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const storyID = Number(id);

  const [activeTab, setActiveTab] = useState<WorkspaceTab>('characters');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // World Data states
  const [story, setStory] = useState<Story | null>(null);
  const [categories, setCategories] = useState<SystemCategory[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [factions, setFactions] = useState<Faction[]>([]);
  const [characters, setCharacters] = useState<Character[]>([]);
  const [plots, setPlots] = useState<Plot[]>([]);
  const [items, setItems] = useState<Item[]>([]);

  // Detected new items for Smart Recheck
  const [detectedRecheckItems, setDetectedRecheckItems] = useState<any[]>([]);
  const [isPublishModalOpen, setIsPublishModalOpen] = useState(false);
  const [isRecheckModalOpen, setIsRecheckModalOpen] = useState(false);

  const fetchAllData = useCallback(async () => {
    if (!storyID) return;
    setLoading(true);
    setError(null);

    try {
      const [
        storyRes,
        catsRes,
        locsRes,
        facsRes,
        charsRes,
        plotsRes,
        itemsRes,
      ] = await Promise.all([
        api.get(`/stories/${storyID}`),
        api.get(`/stories/${storyID}/categories`),
        api.get(`/stories/${storyID}/locations`),
        api.get(`/stories/${storyID}/factions`),
        api.get(`/stories/${storyID}/characters`),
        api.get(`/stories/${storyID}/plots`),
        api.get(`/stories/${storyID}/items`),
      ]);

      setStory(storyRes.data?.data || null);
      setCategories(catsRes.data?.data || []);
      setLocations(locsRes.data?.data || []);
      setFactions(facsRes.data?.data || []);
      setCharacters(charsRes.data?.data || []);
      setPlots(plotsRes.data?.data || []);
      setItems(itemsRes.data?.data || []);
    } catch (err: any) {
      setError(
        err.response?.data?.message ||
          'Không thể tải dữ liệu kịch bản. Vui lòng thử lại.'
      );
    } finally {
      setLoading(false);
    }
  }, [storyID]);

  useEffect(() => {
    fetchAllData();
  }, [fetchAllData]);

  // Map categories sang Record<string, string[]> cho ItemTab
  const categoryMap = categories.reduce((acc, cat) => {
    acc[cat.name] = (cat.items || []).map((i) => i.name);
    return acc;
  }, {} as Record<string, string[]>);

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center text-slate-400">
        <Loader2 className="w-10 h-10 animate-spin text-brand-500 mb-4" />
        <p className="text-sm font-medium">Đang tải không gian sáng tạo kịch bản...</p>
      </div>
    );
  }

  if (error || !story) {
    return (
      <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-6 text-center">
        <div className="w-16 h-16 rounded-2xl bg-rose-500/10 border border-rose-500/20 flex items-center justify-center mb-4 text-rose-400">
          <AlertCircle className="w-8 h-8" />
        </div>
        <h2 className="text-xl font-bold text-white mb-2">Không tìm thấy tác phẩm</h2>
        <p className="text-slate-400 text-sm max-w-md mb-6">{error || 'Tác phẩm không tồn tại hoặc bạn không có quyền truy cập.'}</p>
        <button
          onClick={() => navigate('/dashboard')}
          className="btn-primary flex items-center space-x-2"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Về Bảng Điều Khiển</span>
        </button>
      </div>
    );
  }

  const counts = {
    categories: categories.length,
    locations: locations.length,
    factions: factions.length,
    characters: characters.length,
    plots: plots.length,
    items: items.length,
  };

  return (
    <div className="h-screen flex flex-col bg-slate-950 overflow-hidden text-slate-200">
      {/* Topbar Workspace */}
      <header className="h-14 bg-slate-900 border-b border-slate-800 px-4 flex items-center justify-between flex-shrink-0 z-20">
        <div className="flex items-center space-x-3">
          <button
            onClick={() => navigate('/dashboard')}
            className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition"
            title="Quay lại Bảng Điều Khiển"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>

          <div className="h-5 w-px bg-slate-800" />

          <div className="flex items-center space-x-2.5">
            <BookOpen className="w-5 h-5 text-brand-400" />
            <h1 className="font-bold text-white text-base truncate max-w-xs md:max-w-md">
              {story.title}
            </h1>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-brand-500/20 text-brand-300 border border-brand-500/30">
              {story.genre}
            </span>
            <span
              className={`px-2 py-0.5 rounded text-[11px] font-semibold flex items-center space-x-1 ${
                story.status === 'official'
                  ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'
                  : 'bg-amber-500/20 text-amber-300 border border-amber-500/30'
              }`}
            >
              <CheckCircle2 className="w-3 h-3" />
              <span>{story.status === 'official' ? 'Chính Thức' : 'Bản Nháp'}</span>
            </span>
          </div>
        </div>

        <div className="flex items-center space-x-3">
          {detectedRecheckItems.length > 0 && (
            <button
              onClick={() => setIsRecheckModalOpen(true)}
              className="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl bg-amber-500/20 border border-amber-500/40 text-amber-300 text-xs font-semibold hover:bg-amber-500/30 transition animate-pulse"
              title="Xem các từ khóa mới được phát hiện"
            >
              <Sparkles className="w-3.5 h-3.5" />
              <span>Phát hiện {detectedRecheckItems.length} từ mới</span>
            </button>
          )}

          <button
            onClick={() => setIsPublishModalOpen(true)}
            className="btn-primary text-xs flex items-center space-x-1.5 py-1.5 px-3"
            title="Kiểm tra tính toàn vẹn kịch bản và xuất bản"
          >
            <ShieldCheck className="w-4 h-4" />
            <span>Kiểm Tra Xuất Bản</span>
          </button>

          <button
            onClick={fetchAllData}
            className="p-2 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition"
            title="Làm mới dữ liệu"
          >
            <RefreshCw className="w-4 h-4" />
          </button>
        </div>
      </header>

      {/* Workspace Main Body */}
      <div className="flex flex-1 overflow-hidden">
        {/* Left Sidebar */}
        <WorkspaceSidebar
          activeTab={activeTab}
          onTabChange={setActiveTab}
          counts={counts}
        />

        {/* Center/Right Content Area */}
        <main className="flex-1 bg-slate-950 overflow-hidden relative">
          {activeTab === 'characters' && (
            <CharacterTab
              storyID={storyID}
              characters={characters}
              factions={factions}
              locations={locations}
              categories={categories}
              onRefresh={fetchAllData}
              onDetectedNewItems={(items) => {
                setDetectedRecheckItems(items);
                setIsRecheckModalOpen(true);
              }}
            />
          )}

          {activeTab === 'locations' && (
            <LocationTab
              storyID={storyID}
              locations={locations}
              onRefresh={fetchAllData}
            />
          )}

          {activeTab === 'factions' && (
            <FactionTab
              storyID={storyID}
              factions={factions}
              locations={locations}
              onRefresh={fetchAllData}
            />
          )}

          {activeTab === 'plots' && (
            <PlotTab
              storyID={storyID}
              plots={plots}
              locations={locations}
              characters={characters}
              factions={factions}
              items={items}
              onRefresh={fetchAllData}
            />
          )}

          {activeTab === 'items' && (
            <ItemTab
              storyID={storyID}
              items={items}
              characters={characters}
              factions={factions}
              locations={locations}
              categories={categoryMap}
              onRefresh={fetchAllData}
            />
          )}

          {activeTab === 'categories' && (
            <CategoryTab
              storyID={storyID}
              categories={categories}
              onRefresh={fetchAllData}
            />
          )}
        </main>
      </div>

      {/* Smart Recheck Modal */}
      <SmartRecheckModal
        isOpen={isRecheckModalOpen}
        onClose={() => setIsRecheckModalOpen(false)}
        detectedItems={detectedRecheckItems}
        storyID={storyID}
        categories={categories}
        onSuccess={() => {
          fetchAllData();
          setDetectedRecheckItems([]);
        }}
      />

      {/* Pre-Publish Integrity Verification Modal */}
      <PublishModal
        isOpen={isPublishModalOpen}
        onClose={() => setIsPublishModalOpen(false)}
        storyID={storyID}
        storyTitle={story.title}
        currentStatus={story.status}
        onPublishSuccess={(updatedStory) => {
          setStory(updatedStory);
          fetchAllData();
        }}
        onNavigateTab={(tab) => setActiveTab(tab)}
      />
    </div>
  );
};
