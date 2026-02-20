import { useState, useEffect } from "react";
import { Search } from "lucide-react";
import { useDebounce } from "@/hooks/useDebounce";
import type { ProductFilters } from "@/types";

interface SearchBarProps {
  filters: ProductFilters;
  onFilterChange: (filters: Partial<ProductFilters>) => void;
}

export function SearchBar({ filters, onFilterChange }: SearchBarProps) {
  const [inputValue, setInputValue] = useState(filters.keyword);
  const debouncedValue = useDebounce(inputValue, 300);

  // Sync debounced value to filters
  useEffect(() => {
    if (debouncedValue !== filters.keyword) {
      onFilterChange({ keyword: debouncedValue });
    }
  }, [debouncedValue, onFilterChange, filters.keyword]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const handleSearch = () => {
    onFilterChange({ keyword: inputValue });
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <div className="bg-white px-3 sm:px-4 py-2.5 sm:py-3 border-b border-slate-200 sticky top-[48px] sm:top-[52px] z-40">
      <div className="flex gap-2 sm:gap-3 items-center max-w-7xl mx-auto">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            type="text"
            placeholder="搜索产品..."
            value={inputValue}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            className="w-full pl-9 sm:pl-10 pr-4 py-2.5 sm:py-3 bg-slate-100 rounded-xl text-sm
                       focus:bg-white focus:ring-2 focus:ring-primary-500 focus:outline-none
                       transition-all duration-150 min-h-[44px]"
          />
        </div>
        <button
          onClick={handleSearch}
          className="px-4 sm:px-5 py-2.5 sm:py-3 bg-primary-600 text-white rounded-xl text-sm font-medium
                     hover:bg-primary-700 active:scale-[0.98] transition-all duration-150 cursor-pointer min-h-[44px]
                     flex items-center justify-center"
        >
          <span className="hidden xs:inline">搜索</span>
          <Search className="w-4 h-4 xs:hidden" />
        </button>
      </div>
    </div>
  );
}
