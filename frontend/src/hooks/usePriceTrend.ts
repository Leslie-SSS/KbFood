import { useQuery } from '@tanstack/react-query';
import { productService } from '@/services/productService';

export function usePriceTrend(activityId: string | null, enabled: boolean = false) {
  return useQuery({
    queryKey: ['priceTrend', activityId],
    queryFn: () => productService.getPriceTrend(activityId!),
    enabled: enabled && !!activityId,
    staleTime: 60000, // 1 minute
  });
}
