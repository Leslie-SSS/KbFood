import { UtensilsCrossed, Settings, Ban } from "lucide-react";
import type { SyncStatus } from "@/hooks";

interface NavbarProps {
  productCount: number;
  syncStatus?: SyncStatus;
  onOpenSettings: () => void;
  onOpenBlocked: () => void;
}

// Format last run time: "2026-02-19 23:55:00" -> "02-19 23:55"
function formatLastRunTime(timeStr: string): string {
  const parts = timeStr.split(" ");
  if (parts.length === 2) {
    const date = parts[0].substring(5); // 02-19
    const time = parts[1].substring(0, 5); // 23:55
    return `${date} ${time}`;
  }
  return timeStr;
}

export function Navbar({
  productCount,
  syncStatus,
  onOpenSettings,
  onOpenBlocked,
}: NavbarProps) {
  return (
    <nav className="sticky top-0 z-50 bg-gradient-to-r from-primary-700 to-primary-600 text-white shadow-md">
      <div className="flex items-center justify-between max-w-7xl mx-auto px-3 sm:px-4 py-2 sm:py-2.5">
        <div className="flex items-center gap-2 sm:gap-2.5">
          <UtensilsCrossed className="w-5 h-5 opacity-90" />
          <h1 className="text-sm sm:text-base font-bold tracking-tight">
            美食监控
          </h1>
        </div>
        <div className="flex items-center gap-2 sm:gap-3">
          {/* System status indicator - hidden on very small screens */}
          {syncStatus && (
            <div
              className="hidden sm:flex items-center gap-1.5 text-xs"
              title={
                syncStatus.lastRunTime
                  ? `最后同步: ${syncStatus.lastRunTime}`
                  : "等待同步"
              }
            >
              <span
                className={`w-2 h-2 rounded-full ${
                  syncStatus.isHealthy
                    ? "bg-green-400"
                    : syncStatus.status === "pending"
                      ? "bg-gray-400"
                      : "bg-yellow-400 animate-pulse"
                }`}
              />
              <span className="opacity-80">
                {syncStatus.status === "success"
                  ? "正常"
                  : syncStatus.status === "pending"
                    ? "等待"
                    : syncStatus.status === "running"
                      ? "同步中"
                      : "异常"}
              </span>
            </div>
          )}
          <span className="text-xs opacity-80">{productCount} 商品</span>
          {/* Update time - hidden on small screens */}
          {syncStatus?.lastRunTime && (
            <span className="text-xs opacity-60 hidden md:inline">
              更新于 {formatLastRunTime(syncStatus.lastRunTime)}
            </span>
          )}
          <div className="flex items-center gap-0.5 sm:gap-1">
            <button
              onClick={onOpenSettings}
              className="p-2 sm:p-2 hover:bg-white/10 rounded-lg transition-colors cursor-pointer min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0 flex items-center justify-center"
              title="设置"
            >
              <Settings className="w-4 h-4 sm:w-4 sm:h-4" />
            </button>
            <button
              onClick={onOpenBlocked}
              className="p-2 sm:p-2 hover:bg-white/10 rounded-lg transition-colors cursor-pointer min-w-[44px] min-h-[44px] sm:min-w-0 sm:min-h-0 flex items-center justify-center"
              title="屏蔽列表"
            >
              <Ban className="w-4 h-4 sm:w-4 sm:h-4" />
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}
