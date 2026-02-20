import { useState, useEffect } from "react";
import { settingsService } from "@/services/settingsService";
import { api } from "@/services/api";
import {
  Bell,
  Smartphone,
  Cloud,
  Shield,
  X,
  Send,
  Loader2,
  CheckCircle,
  XCircle,
  ExternalLink,
  ChevronRight,
} from "lucide-react";

interface WelcomeModalProps {
  open: boolean;
  onClose: () => void;
  onComplete: () => void;
}

export function WelcomeModal({ open, onClose, onComplete }: WelcomeModalProps) {
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

  const handleTestNotification = async () => {
    if (!barkKey.trim()) {
      setTestResult({ success: false, message: "请先输入 Bark Key" });
      return;
    }

    setIsTesting(true);
    setTestResult(null);

    try {
      const response = await api.post("/admin/test-notification", {
        barkKey: normalizeBarkKey(barkKey.trim()),
      });

      const data = response.data?.data;

      if (data?.success) {
        setTestResult({
          success: true,
          message: "通知发送成功，请检查您的手机",
        });
      } else {
        const errorMsg = data?.error || response.data?.message || "发送失败";
        setTestResult({
          success: false,
          message: errorMsg,
        });
      }
    } catch (error: unknown) {
      let errorMessage = "请求失败，请检查网络";

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

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const normalizedKey = normalizeBarkKey(barkKey.trim());

      // Save to localStorage
      if (normalizedKey) {
        settingsService.setBarkKey(normalizedKey);
      } else {
        settingsService.clearBarkKey();
      }

      // Save to backend
      try {
        await api.post("/user/settings", { barkKey: normalizedKey });
      } catch (error) {
        console.error("Failed to save settings to backend:", error);
      }

      onComplete();
      onClose();
    } finally {
      setIsSaving(false);
    }
  };

  const handleSkip = () => {
    onClose();
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
        className="relative bg-white rounded-2xl w-full max-w-lg shadow-2xl animate-scale-in overflow-hidden max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="relative bg-gradient-to-br from-primary-600 via-primary-500 to-teal-500 px-5 py-6 sm:px-6 sm:py-7">
          <button
            onClick={onClose}
            className="absolute top-3 right-3 w-8 h-8 flex items-center justify-center rounded-full bg-white/20 hover:bg-white/30 transition-colors cursor-pointer"
          >
            <X className="w-4 h-4 text-white" />
          </button>

          <div className="flex items-center gap-3 pr-8">
            <div className="w-12 h-12 sm:w-14 sm:h-14 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center shadow-lg">
              <Bell className="w-6 h-6 sm:w-7 sm:h-7 text-white" />
            </div>
            <div>
              <h2 className="text-xl sm:text-2xl font-bold text-white">
                欢迎使用美食监控
              </h2>
              <p className="text-sm text-white/80 mt-0.5">
                设置通知，不错过好价
              </p>
            </div>
          </div>
        </div>

        {/* Benefits Section */}
        <div className="p-5 sm:p-6">
          <h3 className="text-base font-semibold text-slate-800 mb-4">
            绑定 Bark 后您可以：
          </h3>

          <div className="space-y-3">
            <BenefitItem
              icon={<Bell className="w-5 h-5" />}
              title="价格提醒推送"
              description="商品降价时，第一时间收到手机推送通知"
              color="text-amber-500"
            />
            <BenefitItem
              icon={<Cloud className="w-5 h-5" />}
              title="数据云同步"
              description="关注商品、屏蔽列表自动同步到云端"
              color="text-blue-500"
            />
            <BenefitItem
              icon={<Smartphone className="w-5 h-5" />}
              title="跨设备访问"
              description="同一账号在手机、平板、电脑上看到相同数据"
              color="text-green-500"
            />
            <BenefitItem
              icon={<Shield className="w-5 h-5" />}
              title="无需注册登录"
              description="Bark Key 即您的唯一标识，简单安全"
              color="text-purple-500"
            />
          </div>

          {/* Divider */}
          <div className="my-5 border-t border-slate-200" />

          {/* Bark Key Input */}
          <div className="space-y-3">
            <label
              htmlFor="barkKey"
              className="flex items-center gap-2 text-sm font-semibold text-slate-700"
            >
              <Bell className="w-4 h-4 text-primary-500" />
              输入您的 Bark Key
            </label>

            <div className="relative">
              <input
                id="barkKey"
                type="text"
                placeholder="完整URL 或 设备Key"
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
              <ExternalLink className="w-4 h-4 text-slate-400 shrink-0 mt-0.5" />
              <div>
                <p className="font-medium text-slate-600">支持格式：</p>
                <p className="text-slate-500">
                  https://api.day.app/XXXXXX 或直接输入 XXXXXX
                </p>
              </div>
            </div>

            {/* Test result */}
            {testResult && (
              <div
                className={`flex items-start gap-3 p-3.5 rounded-xl transition-all duration-300 ${
                  testResult.success
                    ? "bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200"
                    : "bg-gradient-to-r from-red-50 to-orange-50 border border-red-200"
                }`}
              >
                <div
                  className={`w-7 h-7 rounded-full flex items-center justify-center shrink-0 ${
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
        </div>

        {/* Footer */}
        <div className="px-5 py-4 sm:px-6 sm:py-4 bg-slate-50/80 border-t border-slate-100 flex flex-col sm:flex-row gap-2 sm:gap-3">
          <button
            onClick={handleSkip}
            disabled={isSaving || isTesting}
            className="w-full sm:w-auto sm:flex-1 px-4 py-3 bg-white text-slate-600 rounded-xl font-medium
                       border-2 border-slate-200 hover:bg-slate-50 hover:border-slate-300
                       transition-all duration-200 disabled:opacity-50 cursor-pointer text-sm"
          >
            暂不设置
          </button>

          <button
            onClick={handleTestNotification}
            disabled={isTesting || isSaving || !barkKey.trim()}
            className="w-full sm:w-auto sm:flex-1 px-4 py-3 bg-white text-primary-600 rounded-xl font-semibold
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
                测试通知
              </>
            )}
          </button>

          <button
            onClick={handleSave}
            disabled={isSaving || isTesting}
            className="w-full sm:w-auto sm:flex-[1.5] px-4 py-3 bg-gradient-to-r from-primary-600 to-primary-500
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
                完成设置
                <ChevronRight className="w-4 h-4" />
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
            <ExternalLink className="w-3.5 h-3.5" />
            在 App Store 下载 Bark
          </a>
        </div>
      </div>
    </div>
  );
}

// Benefit item component
function BenefitItem({
  icon,
  title,
  description,
  color,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
  color: string;
}) {
  return (
    <div className="flex items-start gap-3 p-3 rounded-xl bg-slate-50/80 hover:bg-slate-50 transition-colors">
      <div className={`w-9 h-9 rounded-lg flex items-center justify-center shrink-0 bg-white shadow-sm ${color}`}>
        {icon}
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-slate-800">{title}</p>
        <p className="text-xs text-slate-500 mt-0.5">{description}</p>
      </div>
    </div>
  );
}

// Normalize Bark key - extract device key from URL if needed
function normalizeBarkKey(input: string): string {
  const trimmed = input.trim();
  if (trimmed.startsWith("http")) {
    const parts = trimmed.split("/");
    return parts[parts.length - 1] || trimmed;
  }
  return trimmed;
}
