import { useMutation, useQueryClient } from '@tanstack/react-query';
import { notificationService } from '@/services/notificationService';
import type { CreateNotificationParams, UpdateNotificationParams } from '@/types';

export function useNotifications() {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: (params: CreateNotificationParams) =>
      notificationService.create(params),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ activityId, params }: { activityId: string; params: UpdateNotificationParams }) =>
      notificationService.update(activityId, params),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (activityId: string) => notificationService.delete(activityId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });

  return {
    create: createMutation.mutateAsync,
    update: updateMutation.mutateAsync,
    delete: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}
