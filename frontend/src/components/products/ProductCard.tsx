import { Bell, BellOff, Ban, TrendingDown, TrendingUp } from "lucide-react";
import type { Product } from "@/types";
import { PriceTrendChart } from "./PriceTrendChart";
import { useInView } from "@/hooks";

interface ProductCardProps {
  product: Product;
  onSetNotification: (product: Product, isEdit: boolean) => void;
  onDeleteNotification: (activityId: string) => void;
  onBlockProduct: (activityId: string) => void;
}

// Format relative time for display
function formatRelativeTime(dateStr: string): string {
  if (!dateStr || dateStr === "0001-01-01T00:00:00Z") return "";

  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return "刚刚";
  if (diffMins < 60) return `${diffMins}分钟前`;
  if (diffHours < 24) return `${diffHours}小时前`;
  if (diffDays < 7) return `${diffDays}天前`;
  return date.toLocaleDateString("zh-CN", { month: "short", day: "numeric" });
}

export function ProductCard({
  product,
  onSetNotification,
  onDeleteNotification,
  onBlockProduct,
}: ProductCardProps) {
  const [chartRef, isInView] = useInView();

  // Calculate price progress (how close to target)
  const priceProgress =
    product.hasNotification && product.targetPrice
      ? Math.min(
          100,
          Math.max(0, (product.targetPrice / product.currentPrice) * 100),
        )
      : null;

  // Is price below target?
  const isBelowTarget =
    product.hasNotification &&
    product.targetPrice &&
    product.currentPrice <= product.targetPrice;

  const handleNotificationClick = () => {
    if (product.hasNotification) {
      onDeleteNotification(product.activityId);
    } else {
      onSetNotification(product, false);
    }
  };

  const handleTargetClick = () => {
    onSetNotification(product, true);
  };

  return (
    <div
      className="bg-white rounded-xl p-3 sm:p-4 border border-slate-200 shadow-card
                 hover:shadow-card-hover hover:-translate-y-0.5
                 transition-all duration-200 flex flex-col gap-2"
    >
      {/* Header */}
      <div className="flex justify-between items-start gap-2">
        <h3 className="text-sm font-semibold text-slate-900 line-clamp-2 leading-tight flex-1">
          {product.title}
        </h3>
        {product.salesStatus !== 1 && (
          <span className="text-[10px] px-2 py-0.5 rounded bg-danger/10 text-danger shrink-0 font-medium">
            已售
          </span>
        )}
      </div>

      {/* Meta info - responsive layout */}
      <div className="flex gap-1.5 text-xs text-slate-500 items-center flex-wrap">
        <span className="shrink-0">{product.platform}</span>
        <span className="text-slate-300">•</span>
        <span className="shrink-0">{product.region}</span>
        <span className="text-slate-300 hidden sm:inline">•</span>
        <span className="truncate max-w-[80px] sm:max-w-[120px] hidden sm:inline">
          {product.shopName}
        </span>
        {product.updateTime && formatRelativeTime(product.updateTime) && (
          <>
            <span className="text-slate-300">•</span>
            <span className="text-slate-400 shrink-0">
              {formatRelativeTime(product.updateTime)}
            </span>
          </>
        )}
      </div>

      {/* Price row */}
      <div className="flex items-baseline gap-2 flex-wrap">
        <span className="text-xl sm:text-2xl font-bold text-danger">
          ¥{product.currentPrice}
        </span>
        {product.originalPrice > 0 &&
          product.originalPrice !== product.currentPrice && (
            <span className="text-xs text-slate-400 line-through">
              ¥{product.originalPrice}
            </span>
          )}
        {product.hasNotification && (
          <span
            className={`text-[10px] px-2 py-0.5 rounded-full ml-auto ${
              isBelowTarget
                ? "bg-success/10 text-success"
                : "bg-primary-100 text-primary-700"
            }`}
          >
            {isBelowTarget ? "已达标" : "监控中"}
          </span>
        )}
      </div>

      {/* Target price info with progress bar */}
      {product.hasNotification && product.targetPrice && (
        <div
          className="mt-1 p-2.5 rounded-lg cursor-pointer transition-colors
                     hover:bg-slate-50 border border-slate-100"
          onClick={handleTargetClick}
        >
          <div className="flex justify-between items-center text-xs mb-1.5">
            <span className="text-slate-500">
              目标价:{" "}
              <span className="font-medium text-primary-600">
                ¥{product.targetPrice}
              </span>
            </span>
            <span className="text-slate-400 flex items-center gap-1">
              {isBelowTarget ? (
                <>
                  <TrendingDown className="w-3 h-3 text-success" />
                  <span className="text-success">已低于目标</span>
                </>
              ) : (
                <>
                  <TrendingUp className="w-3 h-3" />
                  <span>
                    差 ¥
                    {(product.currentPrice - product.targetPrice).toFixed(2)}
                  </span>
                </>
              )}
            </span>
          </div>
          {/* Progress bar */}
          <div className="h-1.5 bg-slate-100 rounded-full overflow-hidden">
            <div
              className={`h-full transition-all duration-500 rounded-full ${
                isBelowTarget
                  ? "bg-gradient-to-r from-success to-emerald-400"
                  : "bg-gradient-to-r from-primary-500 to-primary-400"
              }`}
              style={{ width: `${priceProgress}%` }}
            />
          </div>
        </div>
      )}

      {/* Action buttons - improved touch targets */}
      <div className="flex gap-2 mt-1">
        {product.hasNotification ? (
          <button
            onClick={handleNotificationClick}
            className="flex-1 min-h-[44px] px-3 py-2.5 sm:py-2 text-xs sm:text-sm rounded-lg border border-slate-200 text-slate-500
                       hover:border-slate-300 hover:bg-slate-50 active:scale-[0.98] transition-all duration-150
                       flex items-center justify-center gap-1.5 cursor-pointer"
          >
            <BellOff className="w-4 h-4" />
            <span className="hidden xs:inline">取消监控</span>
            <span className="xs:hidden">取消</span>
          </button>
        ) : (
          <button
            onClick={handleNotificationClick}
            className="flex-1 min-h-[44px] px-3 py-2.5 sm:py-2 text-xs sm:text-sm rounded-lg border border-primary-500 text-primary-600
                       hover:bg-primary-50 active:scale-[0.98] transition-all duration-150
                       flex items-center justify-center gap-1.5 cursor-pointer"
          >
            <Bell className="w-4 h-4" />
            <span>提醒</span>
          </button>
        )}

        <button
          onClick={() => onBlockProduct(product.activityId)}
          className="flex-1 min-h-[44px] px-3 py-2.5 sm:py-2 text-xs sm:text-sm rounded-lg border border-slate-200 text-slate-500
                     hover:border-slate-300 hover:bg-slate-50 active:scale-[0.98] transition-all duration-150
                     flex items-center justify-center gap-1.5 cursor-pointer"
        >
          <Ban className="w-4 h-4" />
          <span>屏蔽</span>
        </button>
      </div>

      {/* Price trend chart - always visible with lazy loading */}
      <div ref={chartRef} className="mt-2 pt-3 border-t border-slate-100">
        {isInView ? (
          <PriceTrendChart
            activityId={product.activityId}
            currentPrice={product.currentPrice}
          />
        ) : (
          <div className="h-20 flex items-center justify-center text-slate-300 text-xs">
            加载中...
          </div>
        )}
      </div>
    </div>
  );
}
