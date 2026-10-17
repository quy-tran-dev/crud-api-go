import React, { useState } from 'react';
import { api } from '../../services/api';
import type { WorkspaceTab } from './WorkspaceSidebar';
import {
  ShieldCheck,
  AlertTriangle,
  CheckCircle2,
  X,
  Loader2,
  RefreshCw,
  ExternalLink,
  Users,
  Film,
  Sword,
  MapPin,
  Shield,
} from 'lucide-react';

interface PublishModalProps {
  isOpen: boolean;
  onClose: () => void;
  storyID: number;
  storyTitle: string;
  currentStatus: 'draft' | 'official';
  onPublishSuccess: (updatedStory: any) => void;
  onNavigateTab: (tab: WorkspaceTab) => void;
}

export const PublishModal: React.FC<PublishModalProps> = ({
  isOpen,
  onClose,
  storyID,
  storyTitle,
  currentStatus,
  onPublishSuccess,
  onNavigateTab,
}) => {
  const [checking, setChecking] = useState(false);
  const [resultState, setResultState] = useState<'idle' | 'failed' | 'success'>('idle');
  const [issues, setIssues] = useState<string[]>([]);
  const [message, setMessage] = useState('');

  if (!isOpen) return null;

  const handleRunVerification = async () => {
    setChecking(true);
    setResultState('idle');
    setIssues([]);
    setMessage('');

    try {
      const response = await api.post(`/stories/${storyID}/publish`);
      setResultState('success');
      setMessage(response.data?.message || 'Xác thực kịch bản và xuất bản chính thức thành công!');
      if (response.data?.data) {
        onPublishSuccess(response.data.data);
      }
    } catch (err: any) {
      setResultState('failed');
      const data = err.response?.data;
      if (data?.issues && Array.isArray(data.issues)) {
        setIssues(data.issues);
        setMessage(data.message || 'Kịch bản chưa đạt tính toàn vẹn để xuất bản');
      } else {
        setMessage(data?.message || 'Có lỗi xảy ra khi xác thực kịch bản.');
      }
    } finally {
      setChecking(false);
    }
  };

  const getSuggestedTab = (issueText: string): WorkspaceTab => {
    const text = issueText.toLowerCase();
    if (text.includes('nhân vật') || text.includes('character')) return 'characters';
    if (text.includes('tình tiết') || text.includes('plot')) return 'plots';
    if (text.includes('thế lực') || text.includes('faction') || text.includes('môn phái')) return 'factions';
    if (text.includes('vật phẩm') || text.includes('item') || text.includes('pháp bảo')) return 'items';
    if (text.includes('địa danh') || text.includes('location') || text.includes('nơi sinh')) return 'locations';
    return 'characters';
  };

  const getTabIcon = (tab: WorkspaceTab) => {
    switch (tab) {
      case 'characters':
        return <Users className="w-3.5 h-3.5 text-amber-400" />;
      case 'locations':
        return <MapPin className="w-3.5 h-3.5 text-emerald-400" />;
      case 'factions':
        return <Shield className="w-3.5 h-3.5 text-blue-400" />;
      case 'plots':
        return <Film className="w-3.5 h-3.5 text-purple-400" />;
      case 'items':
        return <Sword className="w-3.5 h-3.5 text-rose-400" />;
      default:
        return null;
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-xl p-6 shadow-2xl relative">
        {/* Header */}
        <div className="flex items-start justify-between pb-4 border-b border-slate-800">
          <div className="flex items-center space-x-3">
            <div className="p-2.5 rounded-xl bg-brand-500/10 border border-brand-500/20 text-brand-400">
              <ShieldCheck className="w-6 h-6" />
            </div>
            <div>
              <h3 className="text-base font-bold text-white flex items-center space-x-2">
                <span>Pre-Publish Integrity Validator</span>
              </h3>
              <p className="text-xs text-slate-400 mt-0.5">
                Kiểm định toàn vẹn kịch bản trước khi chuyển trạng thái sang Chính Thức
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

        {/* Story Info Banner */}
        <div className="mt-4 p-3 bg-slate-950/60 rounded-xl border border-slate-800 flex items-center justify-between text-xs">
          <div>
            <span className="text-slate-400">Tác phẩm: </span>
            <span className="font-bold text-white ml-1">{storyTitle}</span>
          </div>
          <div className="flex items-center space-x-2">
            <span className="text-slate-400">Hiện tại: </span>
            <span
              className={`px-2 py-0.5 rounded text-[11px] font-semibold ${
                currentStatus === 'official'
                  ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'
                  : 'bg-amber-500/20 text-amber-300 border border-amber-500/30'
              }`}
            >
              {currentStatus === 'official' ? 'Chính Thức' : 'Bản Nháp (Draft)'}
            </span>
          </div>
        </div>

        {/* Body content based on state */}
        <div className="mt-4 min-h-[140px] flex flex-col justify-center">
          {checking && (
            <div className="text-center py-8">
              <Loader2 className="w-10 h-10 animate-spin text-brand-500 mx-auto mb-3" />
              <p className="text-sm font-semibold text-slate-200">
                Đang rà soát tính toàn vẹn hệ thống...
              </p>
              <p className="text-xs text-slate-500 mt-1">
                Kiểm tra liên kết chéo giữa Nhân vật, Thế lực, Bản đồ, Tình tiết và Pháp bảo
              </p>
            </div>
          )}

          {!checking && resultState === 'idle' && (
            <div className="text-center py-6 px-4">
              <p className="text-sm text-slate-300 leading-relaxed">
                Hệ thống sẽ chạy thuật toán quét tự động để đảm bảo mọi tham chiếu ID trong kịch bản (nhân vật, môn phái, địa danh, vật phẩm) đều tồn tại và hợp lệ.
              </p>
              <p className="text-xs text-slate-500 mt-2">
                Nhấn nút bên dưới để bắt đầu quy trình kiểm tra toàn vẹn.
              </p>
            </div>
          )}

          {!checking && resultState === 'success' && (
            <div className="text-center py-6 px-4 bg-emerald-500/5 rounded-xl border border-emerald-500/20">
              <CheckCircle2 className="w-12 h-12 text-emerald-400 mx-auto mb-3" />
              <h4 className="text-base font-bold text-emerald-300">
                Kịch bản đạt chuẩn toàn vẹn 100%!
              </h4>
              <p className="text-xs text-slate-300 mt-2 leading-relaxed">
                {message || 'Tác phẩm đã chính thức được chuyển sang trạng thái Official và sẵn sàng xuất bản.'}
              </p>
            </div>
          )}

          {!checking && resultState === 'failed' && (
            <div className="space-y-3">
              <div className="flex items-center space-x-2 text-rose-400 text-xs font-semibold px-1">
                <AlertTriangle className="w-4 h-4 flex-shrink-0" />
                <span>{message} (Phát hiện {issues.length} điểm cần xử lý):</span>
              </div>

              <div className="max-h-56 overflow-y-auto space-y-2 pr-1">
                {issues.map((issue, idx) => {
                  const targetTab = getSuggestedTab(issue);
                  return (
                    <div
                      key={idx}
                      className="bg-rose-500/5 border border-rose-500/20 rounded-xl p-3 flex items-start justify-between gap-3 text-xs"
                    >
                      <div className="space-y-1">
                        <span className="text-rose-200 leading-relaxed block">
                          • {issue}
                        </span>
                      </div>

                      <button
                        onClick={() => {
                          onClose();
                          onNavigateTab(targetTab);
                        }}
                        className="flex-shrink-0 flex items-center space-x-1 px-2.5 py-1 rounded-lg bg-slate-800 text-slate-200 hover:text-white hover:bg-slate-700 transition text-[11px]"
                        title="Đi đến tab để sửa"
                      >
                        {getTabIcon(targetTab)}
                        <span>Sửa ngay</span>
                        <ExternalLink className="w-3 h-3 text-slate-400 ml-0.5" />
                      </button>
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>

        {/* Footer Actions */}
        <div className="mt-6 pt-4 border-t border-slate-800 flex items-center justify-between">
          <button
            onClick={onClose}
            className="px-4 py-2 text-xs font-medium text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition"
          >
            Đóng
          </button>

          <div className="flex items-center space-x-3">
            {resultState === 'failed' && (
              <button
                onClick={handleRunVerification}
                disabled={checking}
                className="btn-secondary text-xs flex items-center space-x-1.5"
              >
                <RefreshCw className="w-3.5 h-3.5" />
                <span>Kiểm Tra Lại</span>
              </button>
            )}

            {resultState !== 'success' && (
              <button
                onClick={handleRunVerification}
                disabled={checking}
                className="btn-primary text-xs flex items-center space-x-1.5"
              >
                <ShieldCheck className="w-4 h-4" />
                <span>{resultState === 'idle' ? 'Bắt Đầu Kiểm Định' : 'Xuất Bản Lại'}</span>
              </button>
            )}

            {resultState === 'success' && (
              <button
                onClick={onClose}
                className="btn-primary text-xs flex items-center space-x-1.5 bg-emerald-600 hover:bg-emerald-700"
              >
                <CheckCircle2 className="w-4 h-4" />
                <span>Hoàn Tất</span>
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
