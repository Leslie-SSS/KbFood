import { usePriceTrend } from "@/hooks/usePriceTrend";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
} from "chart.js";
import { Line } from "react-chartjs-2";
import { TrendingUp, TrendingDown } from "lucide-react";
import { useMemo } from "react";

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
);

interface PriceTrendChartProps {
  activityId: string;
  currentPrice: number;
}

interface TrendData {
  date: string;
  price: number;
}

// Helper to check if date is today (memoized date at module level)
const getTodayDateString = () => {
  const today = new Date();
  return `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, "0")}-${String(today.getDate()).padStart(2, "0")}`;
};

// Cache today's date string to avoid recalculation
let cachedTodayDate: string | null = null;

function isToday(dateStr: string): boolean {
  if (!cachedTodayDate) {
    cachedTodayDate = getTodayDateString();
  }
  // Compare only the date part (YYYY-MM-DD)
  return dateStr.substring(0, 10) === cachedTodayDate;
}

export function PriceTrendChart({
  activityId,
  currentPrice,
}: PriceTrendChartProps) {
  const { data: trends, isLoading, error } = usePriceTrend(activityId, true);

  // Deduplicate trends by date (keep the lowest price for each date)
  const dedupedTrends: TrendData[] = useMemo(() => {
    if (!trends || trends.length === 0) return [];

    const dateMap = new Map<string, number>();
    for (const t of trends) {
      const dateKey = t.date.substring(0, 10); // Ensure we use date only
      const existing = dateMap.get(dateKey);
      if (existing === undefined || t.price < existing) {
        dateMap.set(dateKey, t.price);
      }
    }

    // Convert back to array and sort by date
    return Array.from(dateMap.entries())
      .map(([date, price]) => ({ date, price }))
      .sort((a, b) => a.date.localeCompare(b.date));
  }, [trends]);

  // Replace today's price with current product price and create chart data
  const chartTrends: TrendData[] = useMemo(() => {
    if (dedupedTrends.length === 0) return [];

    const todayStr = getTodayDateString();
    return dedupedTrends.map((t) => {
      if (t.date === todayStr) {
        return { ...t, price: currentPrice };
      }
      return t;
    });
  }, [dedupedTrends, currentPrice]);

  // Loading state
  if (isLoading) {
    return (
      <div className="h-24 flex items-center justify-center">
        <div className="flex items-center gap-2 text-slate-400 text-xs">
          <div className="w-4 h-4 border-2 border-slate-300 border-t-primary-500 rounded-full animate-spin" />
          加载趋势...
        </div>
      </div>
    );
  }

  // Empty/error state
  if (error || !trends || trends.length === 0 || dedupedTrends.length === 0) {
    return (
      <div className="h-20 flex flex-col items-center justify-center text-slate-400">
        <TrendingUp className="w-5 h-5 mb-1 opacity-50" />
        <span className="text-xs">暂无价格趋势数据</span>
      </div>
    );
  }

  // Calculate min/max prices from chart data (includes current price for today)
  const prices = chartTrends.map((t) => t.price);
  const minPrice = Math.min(...prices);
  const maxPrice = Math.max(...prices, currentPrice); // Include current price for max
  const latestTrendPrice = prices[prices.length - 1];
  const priceChange = prices.length > 1 ? latestTrendPrice - prices[0] : 0;
  const isUp = priceChange > 0;

  // Calculate today's high price and discount (use original dedupedTrends for today's high)
  const todayTrends = dedupedTrends.filter((t) => isToday(t.date));
  const todayHighPrice =
    todayTrends.length > 0
      ? Math.max(...todayTrends.map((t) => t.price))
      : null;
  const todayDiscount =
    todayHighPrice !== null && todayHighPrice > currentPrice
      ? ((todayHighPrice - currentPrice) / todayHighPrice) * 100
      : null;

  const chartData = {
    labels: chartTrends.map((t) => t.date.substring(5)), // MM-DD format
    datasets: [
      {
        label: "价格",
        data: chartTrends.map((t) => t.price),
        borderColor: "#0f766e",
        backgroundColor: (context: { chart: ChartJS }) => {
          const ctx = context.chart.ctx;
          const gradient = ctx.createLinearGradient(0, 0, 0, 120);
          gradient.addColorStop(0, "rgba(15, 118, 110, 0.15)");
          gradient.addColorStop(1, "rgba(15, 118, 110, 0)");
          return gradient;
        },
        borderWidth: 2,
        pointRadius: 3,
        pointBackgroundColor: "#0f766e",
        pointHoverRadius: 5,
        tension: 0.3,
        fill: true,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: "rgba(15, 23, 42, 0.9)",
        titleColor: "#fff",
        bodyColor: "#fff",
        padding: 8,
        cornerRadius: 8,
        displayColors: false,
        callbacks: {
          label: (context: { raw: unknown }) => {
            const value = typeof context.raw === "number" ? context.raw : 0;
            return `¥${value.toFixed(2)}`;
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: false,
        grid: {
          color: "#f1f5f9",
        },
        ticks: {
          color: "#94a3b8",
          font: { size: 10 },
          callback: (value: string | number) => `¥${Number(value).toFixed(2)}`,
        },
      },
      x: {
        grid: {
          display: false,
        },
        ticks: {
          color: "#94a3b8",
          font: { size: 10 },
        },
      },
    },
  };

  return (
    <div className="animate-fade-in">
      {/* Price stats */}
      <div className="flex items-center justify-between text-xs mb-2 px-1">
        <div className="flex items-center gap-3">
          <span className="text-slate-400">
            最低{" "}
            <span className="text-success font-medium">
              ¥{minPrice.toFixed(2)}
            </span>
          </span>
          <span className="text-slate-400">
            最高{" "}
            <span className="text-danger font-medium">
              ¥{maxPrice.toFixed(2)}
            </span>
          </span>
        </div>
        {priceChange !== 0 && (
          <span
            className={`flex items-center gap-0.5 ${isUp ? "text-danger" : "text-success"}`}
          >
            {isUp ? (
              <TrendingUp className="w-3 h-3" />
            ) : (
              <TrendingDown className="w-3 h-3" />
            )}
            {isUp ? "+" : ""}
            {priceChange.toFixed(2)}
          </span>
        )}
      </div>

      {/* Today's discount display */}
      {todayDiscount !== null && todayHighPrice !== null && (
        <div className="text-xs mb-2 px-1 text-success font-medium">
          今日优惠 ↓{todayDiscount.toFixed(1)}% (最高¥
          {todayHighPrice.toFixed(2)})
        </div>
      )}

      {/* Chart */}
      <div className="h-24">
        <Line data={chartData} options={options} />
      </div>
    </div>
  );
}
