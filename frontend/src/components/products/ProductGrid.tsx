import type { Product } from '@/types';
import { ProductCard } from './ProductCard';
import { Search, Package } from 'lucide-react';

interface ProductGridProps {
  products: Product[];
  isLoading: boolean;
  onSetNotification: (product: Product, isEdit: boolean) => void;
  onDeleteNotification: (activityId: string) => void;
  onBlockProduct: (activityId: string) => void;
}

// Skeleton component for loading state
function ProductSkeleton() {
  return (
    <div className="bg-white rounded-xl p-4 border border-slate-200">
      {/* Title skeleton */}
      <div className="h-4 bg-slate-200 rounded w-3/4 mb-3 animate-pulse" />

      {/* Meta skeleton */}
      <div className="flex gap-2 mb-3">
        <div className="h-3 bg-slate-100 rounded w-16 animate-pulse" />
        <div className="h-3 bg-slate-100 rounded w-12 animate-pulse" />
      </div>

      {/* Price skeleton */}
      <div className="h-6 bg-slate-200 rounded w-20 mb-4 animate-pulse" />

      {/* Progress bar skeleton */}
      <div className="h-8 bg-slate-50 rounded mb-4 animate-pulse" />

      {/* Buttons skeleton */}
      <div className="flex gap-2">
        <div className="h-8 bg-slate-100 rounded flex-1 animate-pulse" />
        <div className="h-8 bg-slate-100 rounded flex-1 animate-pulse" />
        <div className="h-8 bg-slate-100 rounded flex-1 animate-pulse" />
      </div>

      {/* Chart skeleton */}
      <div className="mt-3 pt-3 border-t border-slate-100">
        <div className="h-20 bg-slate-50 rounded animate-pulse" />
      </div>
    </div>
  );
}

export function ProductGrid({
  products,
  isLoading,
  onSetNotification,
  onDeleteNotification,
  onBlockProduct,
}: ProductGridProps) {
  // Loading state with skeleton
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <ProductSkeleton key={i} />
        ))}
      </div>
    );
  }

  // Empty state
  if (products.length === 0) {
    return (
      <div className="text-center py-20">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-100 mb-4">
          <Search className="w-8 h-8 text-slate-400" />
        </div>
        <h3 className="text-lg font-medium text-slate-700 mb-2">未找到相关产品</h3>
        <p className="text-sm text-slate-500">尝试调整筛选条件或搜索关键词</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {/* Results count */}
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Package className="w-4 h-4" />
        <span>共 {products.length} 个产品</span>
      </div>

      {/* Product grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {products.map((product) => (
          <ProductCard
            key={product.activityId}
            product={product}
            onSetNotification={onSetNotification}
            onDeleteNotification={onDeleteNotification}
            onBlockProduct={onBlockProduct}
          />
        ))}
      </div>
    </div>
  );
}
