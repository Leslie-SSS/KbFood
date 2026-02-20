import { useState, useEffect } from "react";
import { settingsService } from "@/services/settingsService";
import { api } from "@/services/api";
import {
  Bell,
  Send,
  Loader2,
  CheckCircle,
  XCircle,
  X,
  ExternalLink,
  HelpCircle,
  Smartphone,
} from "lucide-react";

interface SettingsModalProps {
  open: boolean;
  onClose: () => void;
  onSave: (message: string) => void;
}

export function SettingsModal({ open, onClose, onSave }: SettingsModalProps) {
  const [barkKey, setBarkKey] = useState("");
  const [isTesting, setIsTesting] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    message: string;
  } | null>(null);

  useEffect(() => {
    if (open) {
      setBarkKey(settingsService.getBarkKey());
      setTestResult(null);
    }
  }, [open]);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const normalizedKey = barkKey.trim();

      // 1. Save to localStorage (for frontend use)
      if (normalizedKey) {
        settingsService.setBarkKey(normalizedKey);
      } else {
        settingsService.clearBarkKey();
      }

      // 2. Save to backend (for background notifications)
      try {
        await api.post("/user/settings", { barkKey: normalizedKey });
      } catch (error) {
        console.error("Failed to save settings to backend:", error);
        // Continue even if backend save fails - localStorage is the primary source
      }

      onSave(normalizedKey ? "设置已保存" : "已清除 Bark Key");
      onClose();
    } finally {
      setIsSaving(false);
    }
  };

  const handleTestNotification = async () => {
    if (!barkKey.trim()) {
      setTestResult({ success: false, message: "请先输入 Bark Key" });
      return;
    }

    setIsTesting(true);
    setTestResult(null);

    try {
      const response = await api.post("/admin/test-notification", {
        barkKey: barkKey.trim(),
      });

      const data = response.data?.data;

      if (data?.success) {
        setTestResult({
          success: true,
          message: "通知发送成功，请检查您的手机",
        });
      } else {
        // Get detailed error from server response
        const errorMsg = data?.error || response.data?.message || "发送失败";
        setTestResult({
          success: false,
          message: errorMsg,
        });
      }
    } catch (error: unknown) {
      // Detailed error handling
      let errorMessage = "请求失败";

      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as {
          response?: { data?: { message?: string }; status?: number };
        };
        if (axiosError.response?.data?.message) {
          errorMessage = axiosError.response.data.message;
        } else if (axiosError.response?.status) {
          errorMessage = `请求失败 (${axiosError.response.status})`;
        }
      } else if (error && typeof error === "object" && "message" in error) {
        errorMessage = String((error as { message: string }).message);
      }

      setTestResult({
        success: false,
        message: errorMessage,
      });
    } finally {
      setIsTesting(false);
    }
  };

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      onClick={onClose}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" />

      {/* Modal */}
      <div
        className="relative bg-white rounded-2xl w-full max-w-md shadow-2xl animate-scale-in overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="relative bg-gradient-to-br from-primary-600 via-primary-500 to-teal-500 px-5 py-5">
          <button
            onClick={onClose}
            className="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full bg-white/20 hover:bg-white/30 transition-colors cursor-pointer"
          >
            <X className="w-4 h-4 text-white" />
          </button>

          <div className="flex items-center gap-3 pr-8">
            <div className="w-11 h-11 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center shadow-lg">
              <Bell className="w-5 h-5 text-white" />
            </div>
            <div>
              <h2 className="text-lg font-bold text-white">通知设置</h2>
              <p className="text-xs text-white/70">配置价格提醒推送</p>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="p-5 space-y-5">
          {/* Bark Key Input */}
          <div className="space-y-2">
            <label
              htmlFor="barkKey"
              className="flex items-center gap-1.5 text-sm font-semibold text-slate-700"
            >
              <Smartphone className="w-4 h-4 text-primary-500" />
              Bark 推送地址
            </label>

            <div className="relative">
              <input
                id="barkKey"
                type="text"
                placeholder="输入完整URL或设备Key"
                value={barkKey}
                onChange={(e) => {
                  setBarkKey(e.target.value);
                  setTestResult(null);
                }}
                className="w-full px-4 py-3.5 border-2 rounded-xl text-sm
                           transition-all duration-200 outline-none placeholder:text-slate-300
                           focus:border-primary-400 focus:bg-primary-50/30
                           border-slate-200 bg-white hover:border-slate-300"
              />
              {barkKey.trim() && (
                <button
                  onClick={() => {
                    setBarkKey("");
                    setTestResult(null);
                  }}
                  className="absolute right-3 top-1/2 -translate-y-1/2 w-6 h-6 flex items-center justify-center rounded-full bg-slate-100 hover:bg-slate-200 transition-colors cursor-pointer"
                >
                  <X className="w-3.5 h-3.5 text-slate-400" />
                </button>
              )}
            </div>

            {/* Format hint */}
            <div className="flex items-start gap-2 text-xs text-slate-500 bg-slate-50 rounded-lg px-3 py-2.5">
              <HelpCircle className="w-4 h-4 text-slate-400 shrink-0 mt-0.5" />
              <div className="space-y-1">
                <p className="font-medium text-slate-600">支持两种格式：</p>
                <p className="text-slate-500">
                  • 完整地址：https://api.day.app/XXXXXX
                </p>
                <p className="text-slate-500">• 设备 Key：XXXXXX</p>
              </div>
            </div>
          </div>

          {/* Test result feedback */}
          {testResult && (
            <div
              className={`flex items-start gap-3 p-4 rounded-xl transition-all duration-300 ${
                testResult.success
                  ? "bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200"
                  : "bg-gradient-to-r from-red-50 to-orange-50 border border-red-200"
              }`}
            >
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 ${
                  testResult.success ? "bg-green-100" : "bg-red-100"
                }`}
              >
                {testResult.success ? (
                  <CheckCircle className="w-4 h-4 text-green-600" />
                ) : (
                  <XCircle className="w-4 h-4 text-red-600" />
                )}
              </div>
              <div className="flex-1 min-w-0">
                <p
                  className={`text-sm font-medium ${
                    testResult.success ? "text-green-700" : "text-red-700"
                  }`}
                >
                  {testResult.success ? "发送成功" : "发送失败"}
                </p>
                <p
                  className={`text-xs mt-0.5 ${
                    testResult.success ? "text-green-600" : "text-red-600"
                  }`}
                >
                  {testResult.message}
                </p>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-5 py-4 bg-slate-50/80 border-t border-slate-100 flex gap-3">
          <button
            onClick={onClose}
            disabled={isSaving || isTesting}
            className="flex-1 px-4 py-3 bg-white text-slate-600 rounded-xl font-medium
                       border-2 border-slate-200 hover:bg-slate-50 hover:border-slate-300
                       transition-all duration-200 disabled:opacity-50 cursor-pointer text-sm"
          >
            取消
          </button>

          <button
            onClick={handleTestNotification}
            disabled={isTesting || isSaving || !barkKey.trim()}
            className="flex-1 px-4 py-3 bg-white text-primary-600 rounded-xl font-semibold
                       border-2 border-primary-200 hover:bg-primary-50 hover:border-primary-300
                       disabled:opacity-40 disabled:cursor-not-allowed
                       transition-all duration-200
                       flex items-center justify-center gap-2 cursor-pointer text-sm"
          >
            {isTesting ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                发送中
              </>
            ) : (
              <>
                <Send className="w-4 h-4" />
                测试
              </>
            )}
          </button>

          <button
            onClick={handleSave}
            disabled={isSaving || isTesting}
            className="flex-[1.5] px-4 py-3 bg-gradient-to-r from-primary-600 to-primary-500
                       text-white rounded-xl font-semibold shadow-lg shadow-primary-500/20
                       hover:from-primary-700 hover:to-primary-600 hover:shadow-xl
                       disabled:opacity-40 disabled:cursor-not-allowed disabled:shadow-none
                       active:scale-[0.98] transition-all duration-200
                       flex items-center justify-center gap-2 cursor-pointer text-sm"
          >
            {isSaving ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                保存中
              </>
            ) : (
              <>
                <Bell className="w-4 h-4" />
                保存设置
              </>
            )}
          </button>
        </div>

        {/* Help link */}
        <div className="px-5 py-3 bg-slate-50 border-t border-slate-100">
          <a
            href="https://apps.apple.com/app/bark/id1403753865"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center justify-center gap-1.5 text-xs text-slate-400 hover:text-primary-500 transition-colors cursor-pointer"
          >
            <ExternalLink className="w-3.5 h-3.5" />在 App Store 下载 Bark
          </a>
        </div>
      </div>
    </div>
  );
}
