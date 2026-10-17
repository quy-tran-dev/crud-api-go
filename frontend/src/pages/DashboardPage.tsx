import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import type { Story } from '../types';
import { Navbar } from '../components/Navbar';
import {
  BookOpen,
  Plus,
  Search,
  Sparkles,
  Trash2,
  FolderOpen,
  Calendar,
  AlertCircle,
  X,
  Loader2,
  CheckCircle2,
} from 'lucide-react';

export const DashboardPage: React.FC = () => {
  const [stories, setStories] = useState<Story[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'draft' | 'official'>('all');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Form State tạo mới tác phẩm
  const [title, setTitle] = useState('');
  const [genre, setGenre] = useState('Tu Tiên');
  const [synopsis, setSynopsis] = useState('');
  const [worldSetting, setWorldSetting] = useState('');

  const navigate = useNavigate();

  const fetchStories = async () => {
    try {
      setIsLoading(true);
      const params: any = { page: 1, limit: 100 };
      if (statusFilter !== 'all') params.status = statusFilter;
      if (search) params.search = search;

      const response = await api.get('/stories', { params });
      setStories(response.data.data || []);
    } catch (err: any) {
      console.error('Lỗi tải danh sách tác phẩm:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchStories();
  }, [statusFilter]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    fetchStories();
  };

  const handleCreateStory = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      await api.post('/stories', {
        title,
        genre,
        synopsis,
        world_setting: worldSetting,
      });
      setIsModalOpen(false);
      // Reset form
      setTitle('');
      setGenre('Tu Tiên');
      setSynopsis('');
      setWorldSetting('');
      fetchStories();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Không thể tạo tác phẩm. Vui lòng thử lại.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteStory = async (id: number, title: string) => {
    if (!window.confirm(`Bạn có chắc chắn muốn xóa tác phẩm "${title}"? Thao tác này sẽ xóa toàn bộ nhân vật, địa danh, tình tiết và vật phẩm bên trong!`)) {
      return;
    }

    try {
      await api.delete(`/stories/${id}`);
      setStories(stories.filter((s) => s.id !== id));
    } catch (err: any) {
      alert('Lỗi khi xóa tác phẩm: ' + (err.response?.data?.message || err.message));
    }
  };

  const popularGenres = [
    { name: 'Tu Tiên', desc: 'Tự động tạo Cảnh giới, Tông môn, Phẩm cấp đan dược' },
    { name: 'Cyberpunk', desc: 'Tự động tạo Tập đoàn, Cyberware, Băng nhóm' },
    { name: 'Trinh Thám', desc: 'Tự động tạo Manh mối, Cục cảnh sát, Bằng chứng' },
    { name: 'Huyền Huyễn', desc: 'Tự động tạo Hệ thống ma pháp, Chủng tộc, Gia tộc' },
  ];

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
      <Navbar />

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header Dashboard */}
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 pb-8 border-b border-slate-800/80">
          <div>
            <h1 className="text-3xl font-extrabold text-white tracking-tight flex items-center space-x-3">
              <span>Bàn Làm Việc Kịch Bản</span>
              <span className="text-xs font-semibold px-2.5 py-1 bg-brand-950 text-brand-400 border border-brand-800 rounded-full">
                {stories.length} Tác phẩm
              </span>
            </h1>
            <p className="mt-1 text-sm text-slate-400">
              Quản lý các vũ trụ tiểu thuyết, phân cấp thư mục và hệ thống nhân vật của bạn
            </p>
          </div>

          <button
            onClick={() => setIsModalOpen(true)}
            className="inline-flex items-center justify-center space-x-2 bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-medium px-4 py-2.5 rounded-xl shadow-lg shadow-brand-900/40 transition-all text-sm"
          >
            <Plus className="w-5 h-5" />
            <span>Tạo Tác Phẩm Mới</span>
          </button>
        </div>

        {/* Filter & Search Bar */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 my-6">
          <form onSubmit={handleSearchSubmit} className="relative flex-1 max-w-md">
            <Search className="w-4 h-4 absolute left-3.5 top-3 text-slate-500" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Tìm kiếm theo tên truyện hoặc thể loại..."
              className="w-full pl-10 pr-4 py-2 bg-slate-900/80 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-brand-500"
            />
          </form>

          {/* Status Tabs */}
          <div className="flex items-center space-x-2 bg-slate-900/90 p-1 rounded-xl border border-slate-800 text-xs">
            <button
              onClick={() => setStatusFilter('all')}
              className={`px-3 py-1.5 rounded-lg font-medium transition-all ${
                statusFilter === 'all'
                  ? 'bg-brand-600 text-white shadow'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Tất cả
            </button>
            <button
              onClick={() => setStatusFilter('draft')}
              className={`px-3 py-1.5 rounded-lg font-medium transition-all ${
                statusFilter === 'draft'
                  ? 'bg-amber-600 text-white shadow'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Bản Nháp
            </button>
            <button
              onClick={() => setStatusFilter('official')}
              className={`px-3 py-1.5 rounded-lg font-medium transition-all ${
                statusFilter === 'official'
                  ? 'bg-emerald-600 text-white shadow'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Chính Thức
            </button>
          </div>
        </div>

        {/* Story Cards Grid */}
        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-20 text-slate-500 space-y-3">
            <Loader2 className="w-8 h-8 animate-spin text-brand-500" />
            <p className="text-sm">Đang tải danh sách tác phẩm...</p>
          </div>
        ) : stories.length === 0 ? (
          <div className="text-center py-20 border border-dashed border-slate-800 rounded-3xl bg-slate-900/30 p-8">
            <BookOpen className="w-12 h-12 text-slate-600 mx-auto mb-3" />
            <h3 className="text-lg font-semibold text-slate-300">Chưa có tác phẩm nào</h3>
            <p className="text-sm text-slate-500 max-w-sm mx-auto mt-1 mb-5">
              Hãy bắt đầu sáng tạo bằng cách bấm nút "Tạo Tác Phẩm Mới" để hệ thống tự động sinh bộ Từ điển hằng số mẫu.
            </p>
            <button
              onClick={() => setIsModalOpen(true)}
              className="inline-flex items-center space-x-2 bg-slate-800 hover:bg-slate-700 text-brand-300 px-4 py-2 rounded-xl text-sm font-medium border border-slate-700 transition-colors"
            >
              <Plus className="w-4 h-4" />
              <span>Khởi tạo câu chuyện đầu tiên</span>
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {stories.map((story) => (
              <div
                key={story.id}
                className="group bg-slate-900/90 border border-slate-800/90 hover:border-brand-500/50 rounded-2xl p-5 flex flex-col justify-between transition-all hover:shadow-xl hover:shadow-brand-950/40 relative overflow-hidden"
              >
                {/* Background glow decoration */}
                <div className="absolute -top-12 -right-12 w-28 h-28 bg-brand-500/10 rounded-full blur-2xl group-hover:bg-brand-500/20 transition-all pointer-events-none" />

                <div>
                  <div className="flex items-start justify-between gap-3 mb-3">
                    <span className="text-xs font-semibold px-2.5 py-1 rounded-md bg-slate-800 text-slate-300 border border-slate-700/80">
                      {story.genre || 'Chưa phân loại'}
                    </span>

                    {story.is_draft || story.status === 'draft' ? (
                      <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-amber-950/80 text-amber-400 border border-amber-800/60 inline-flex items-center space-x-1">
                        <span className="w-1.5 h-1.5 rounded-full bg-amber-400 animate-pulse" />
                        <span>Bản Nháp</span>
                      </span>
                    ) : (
                      <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-emerald-950/80 text-emerald-400 border border-emerald-800/60 inline-flex items-center space-x-1">
                        <CheckCircle2 className="w-3.5 h-3.5" />
                        <span>Chính Thức</span>
                      </span>
                    )}
                  </div>

                  <h3 className="text-xl font-bold text-white group-hover:text-brand-300 transition-colors line-clamp-1 mb-2">
                    {story.title}
                  </h3>

                  <p className="text-sm text-slate-400 line-clamp-3 mb-4 leading-relaxed">
                    {story.synopsis || 'Chưa có tóm tắt cho câu chuyện này...'}
                  </p>
                </div>

                <div className="pt-4 border-t border-slate-800/80 flex items-center justify-between">
                  <div className="flex items-center space-x-1 text-xs text-slate-500">
                    <Calendar className="w-3.5 h-3.5" />
                    <span>{new Date(story.created_at).toLocaleDateString('vi-VN')}</span>
                  </div>

                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => handleDeleteStory(story.id, story.title)}
                      className="p-2 text-slate-500 hover:text-red-400 hover:bg-slate-800 rounded-lg transition-colors"
                      title="Xóa tác phẩm"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>

                    <button
                      onClick={() => navigate(`/stories/${story.id}`)}
                      className="inline-flex items-center space-x-1.5 bg-brand-600 hover:bg-brand-500 text-white text-xs font-medium px-3 py-1.5 rounded-lg shadow transition-colors"
                    >
                      <FolderOpen className="w-3.5 h-3.5" />
                      <span>Vào Kịch Bản</span>
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>

      {/* MODAL TẠO TÁC PHẨM MỚI */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl max-w-xl w-full p-6 shadow-2xl relative animate-in fade-in zoom-in-95 duration-200">
            <button
              onClick={() => setIsModalOpen(false)}
              className="absolute top-4 right-4 text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition-colors"
            >
              <X className="w-5 h-5" />
            </button>

            <div className="flex items-center space-x-3 mb-4">
              <div className="p-2 bg-brand-950 border border-brand-800/80 rounded-xl text-brand-400">
                <Sparkles className="w-6 h-6" />
              </div>
              <div>
                <h3 className="text-lg font-bold text-white">Khởi Tạo Vũ Trụ Kịch Bản</h3>
                <p className="text-xs text-slate-400">Hệ thống sẽ tự động cấu hình bộ Từ điển theo thể loại bạn chọn</p>
              </div>
            </div>

            {error && (
              <div className="mb-4 flex items-center space-x-2 bg-red-950/60 border border-red-800 p-3 rounded-xl text-red-300 text-xs">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleCreateStory} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tên Tác Phẩm / Tựa Truyện *
                </label>
                <input
                  type="text"
                  required
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Ví dụ: Đại Chúa Tể Càn Khôn, Mật Mã Bóng Đêm..."
                  className="w-full px-3.5 py-2 bg-slate-950 border border-slate-700/80 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Thể Loại Truyện (Chọn để kích hoạt Auto-Seed)
                </label>
                <div className="grid grid-cols-2 gap-2 mb-2">
                  {popularGenres.map((g) => (
                    <button
                      type="button"
                      key={g.name}
                      onClick={() => setGenre(g.name)}
                      className={`p-2.5 rounded-xl border text-left text-xs transition-all ${
                        genre === g.name
                          ? 'border-brand-500 bg-brand-950/60 text-white ring-1 ring-brand-500'
                          : 'border-slate-800 bg-slate-950/60 text-slate-400 hover:border-slate-700'
                      }`}
                    >
                      <div className="font-semibold text-slate-200">{g.name}</div>
                      <div className="text-[10px] text-slate-500 line-clamp-1 mt-0.5">{g.desc}</div>
                    </button>
                  ))}
                </div>
                <input
                  type="text"
                  value={genre}
                  onChange={(e) => setGenre(e.target.value)}
                  placeholder="Hoặc tự gõ thể loại khác..."
                  className="w-full px-3 py-1.5 bg-slate-950 border border-slate-700/80 rounded-lg text-xs text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Tóm Tắt Cốt Truyện (Synopsis)
                </label>
                <textarea
                  rows={3}
                  value={synopsis}
                  onChange={(e) => setSynopsis(e.target.value)}
                  placeholder="Giới thiệu khái quát mở đầu, mục tiêu và biến cố..."
                  className="w-full px-3.5 py-2 bg-slate-950 border border-slate-700/80 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase mb-1">
                  Bối Cảnh Thế Giới (World Setting)
                </label>
                <input
                  type="text"
                  value={worldSetting}
                  onChange={(e) => setWorldSetting(e.target.value)}
                  placeholder="Quy tắc ma pháp, không gian, thời đại..."
                  className="w-full px-3.5 py-2 bg-slate-950 border border-slate-700/80 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-brand-500"
                />
              </div>

              <div className="pt-2 flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 rounded-xl text-sm font-medium text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="inline-flex items-center space-x-2 bg-brand-600 hover:bg-brand-500 text-white text-sm font-semibold px-5 py-2 rounded-xl shadow-lg shadow-brand-900/40 transition-all disabled:opacity-50"
                >
                  {isSubmitting ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      <span>Đang khởi tạo...</span>
                    </>
                  ) : (
                    <span>Tạo Tác Phẩm</span>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
