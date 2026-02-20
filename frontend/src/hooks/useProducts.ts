import { useQuery, useQueryClient } from '@tanstack/react-query';
import { productService } from '@/services/productService';
import type { ProductFilters, Product } from '@/types';

export function useProducts(filters: ProductFilters) {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ['products', filters],
    queryFn: () => productService.getProducts({
      keyword: filters.keyword,
      platform: filters.platform,
      region: filters.region,
      salesStatus: filters.salesStatus,
      monitorStatus: filters.monitorStatus,
      recentSevenDays: filters.recentSevenDays,
    }),
    staleTime: 30000, // 30 seconds
  });

  // Sort products: monitored first, then by price gap
  const sortedProducts: Product[] = query.data
    ? [...query.data].sort((a, b) => {
        const monitorA = a.hasNotification ? 1 : 0;
        const monitorB = b.hasNotification ? 1 : 0;

        if (monitorA !== monitorB) {
          return monitorB - monitorA;
        }

        if (monitorA === 1) {
          const gapA = a.currentPrice - (a.targetPrice || 0);
          const gapB = b.currentPrice - (b.targetPrice || 0);
          return gapA - gapB;
        }

        return 0;
      })
    : [];

  const refreshProducts = () => {
    queryClient.invalidateQueries({ queryKey: ['products'] });
  };

  return {
    products: sortedProducts,
    isLoading: query.isLoading,
    error: query.error,
    refreshProducts,
  };
}
