import {
  PLATFORMS,
  REGIONS,
  SALES_STATUS,
  MONITOR_STATUS,
  type ProductFilters,
} from "@/types";
import { ChevronDown, X } from "lucide-react";
import { useState, useRef, useEffect, useCallback } from "react";
import { createPortal } from "react-dom";

interface ProductFiltersPanelProps {
  filters: ProductFilters;
  onFilterChange: (filters: Partial<ProductFilters>) => void;
  onReset: () => void;
}

// Tag-style dropdown filter component with Portal for proper z-index
function FilterTag({
  label,
  options,
  value,
  onChange,
}: {
  label: string;
  options: readonly { value: string; label: string }[];
  value: string;
  onChange: (value: string) => void;
}) {
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0, width: 0 });
  const ref = useRef<HTMLDivElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const selected = options.find((o) => o.value === value);

  const updatePosition = useCallback(() => {
    if (ref.current) {
      const rect = ref.current.getBoundingClientRect();
      setPosition({
        top: rect.bottom + window.scrollY,
        left: rect.left,
        width: rect.width,
      });
    }
  }, []);

  useEffect(() => {
    if (open) {
      updatePosition();
    }
  }, [open, updatePosition]);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      const target = e.target as Node;
      // Check if click is inside the button or dropdown
      if (ref.current?.contains(target) || dropdownRef.current?.contains(target)) {
        return;
      }
      setOpen(false);
    }

    if (open) {
      // Use pointerdown for better compatibility
      document.addEventListener("pointerdown", handleClickOutside);
      return () => document.removeEventListener("pointerdown", handleClickOutside);
    }
  }, [open]);

  const handleSelect = useCallback((optValue: string) => {
    onChange(optValue);
    setOpen(false);
  }, [onChange]);

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className={`flex items-center gap-1.5 px-3 py-2 sm:py-1.5 rounded-full text-xs sm:text-sm whitespace-nowrap transition-all cursor-pointer min-h-[36px] sm:min-h-0
                   ${
                     value
                       ? "bg-primary-100 text-primary-700 font-medium"
                       : "bg-slate-100 text-slate-600 hover:bg-slate-200"
                   }`}
      >
        <span>{selected?.label || label}</span>
        <ChevronDown
          className={`w-3.5 h-3.5 transition-transform ${open ? "rotate-180" : ""}`}
        />
      </button>
      {open && createPortal(
        <div
          ref={dropdownRef}
          className="fixed bg-white border border-slate-200 rounded-lg shadow-xl py-1 z-[9999]"
          style={{
            top: position.top + 4,
            left: position.left,
            minWidth: Math.max(position.width, 140),
          }}
        >
          {options.map((opt) => (
            <button
              type="button"
              key={opt.value}
              onMouseDown={(e) => {
                e.preventDefault();
                e.stopPropagation();
              }}
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                handleSelect(opt.value);
              }}
              className={`w-full px-3 py-2.5 sm:py-2 text-xs sm:text-sm text-left hover:bg-slate-50 transition-colors cursor-pointer
                         ${opt.value === value ? "text-primary-600 font-medium bg-primary-50" : "text-slate-600"}`}
            >
              {opt.label}
            </button>
          ))}
        </div>,
        document.body
      )}
    </div>
  );
}

export function ProductFilters({
  filters,
  onFilterChange,
  onReset,
}: ProductFiltersPanelProps) {
  // Count active filters
  const activeFilterCount = [
    filters.platform,
    filters.region,
    filters.salesStatus,
    filters.monitorStatus,
    filters.recentSevenDays,
  ].filter(Boolean).length;

  return (
    <div className="py-2 sm:py-2.5 bg-slate-50 border-b border-slate-200">
      <div className="px-3 sm:px-4 max-w-7xl mx-auto">
        {/* Horizontal scrollable container on mobile */}
        <div className="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0 sm:flex-wrap scrollbar-hide">
          <span className="text-xs text-slate-500 shrink-0 hidden sm:inline">
            筛选:
          </span>

          {/* Platform filter */}
          <FilterTag
            label="平台"
            options={PLATFORMS}
            value={filters.platform}
            onChange={(v) => onFilterChange({ platform: v })}
          />

          {/* Region filter */}
          <FilterTag
            label="地区"
            options={REGIONS}
            value={filters.region}
            onChange={(v) => onFilterChange({ region: v })}
          />

          {/* Sales status filter */}
          <FilterTag
            label="状态"
            options={SALES_STATUS}
            value={filters.salesStatus}
            onChange={(v) => onFilterChange({ salesStatus: v })}
          />

          {/* Monitor status filter */}
          <FilterTag
            label="监控"
            options={MONITOR_STATUS}
            value={filters.monitorStatus}
            onChange={(v) => onFilterChange({ monitorStatus: v })}
          />

          {/* Recent 7 days toggle */}
          <button
            type="button"
            onClick={() =>
              onFilterChange({ recentSevenDays: !filters.recentSevenDays })
            }
            className={`px-3 py-2 sm:py-1.5 rounded-full text-xs sm:text-sm whitespace-nowrap transition-all cursor-pointer min-h-[36px] sm:min-h-0
                       ${
                         filters.recentSevenDays
                           ? "bg-primary-100 text-primary-700 font-medium"
                           : "bg-slate-100 text-slate-600 hover:bg-slate-200"
                       }`}
          >
            近7天上新
          </button>

          {/* Reset button */}
          {activeFilterCount > 0 && (
            <button
              type="button"
              onClick={onReset}
              className="flex items-center gap-1 px-3 py-2 sm:py-1.5 rounded-full text-xs sm:text-sm text-slate-500
                         hover:text-primary-600 hover:bg-primary-50 transition-all cursor-pointer shrink-0 min-h-[36px] sm:min-h-0"
            >
              <X className="w-3.5 h-3.5" />
              <span className="hidden xs:inline">重置</span>
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
