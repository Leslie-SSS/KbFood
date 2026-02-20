import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { productService } from '@/services/productService';

export function useBlockedProducts() {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ['blockedProducts'],
    queryFn: () => productService.getBlockedProducts(),
  });

  const blockMutation = useMutation({
    mutationFn: (activityId: string) => productService.blockProduct(activityId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      queryClient.invalidateQueries({ queryKey: ['blockedProducts'] });
    },
  });

  const unblockMutation = useMutation({
    mutationFn: (activityId: string) => productService.unblockProduct(activityId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      queryClient.invalidateQueries({ queryKey: ['blockedProducts'] });
    },
  });

  return {
    blockedProducts: query.data || [],
    isLoading: query.isLoading,
    block: blockMutation.mutateAsync,
    unblock: unblockMutation.mutateAsync,
    isBlocking: blockMutation.isPending,
    isUnblocking: unblockMutation.isPending,
  };
}
