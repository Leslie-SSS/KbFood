import { api } from './api';
import type { ApiResponse, CreateNotificationParams, UpdateNotificationParams } from '@/types';

export const notificationService = {
  // Create a new notification
  create: async (params: CreateNotificationParams): Promise<void> => {
    await api.post<ApiResponse<null>>('/products/notifications', params);
  },

  // Update an existing notification
  update: async (activityId: string, params: UpdateNotificationParams): Promise<void> => {
    await api.put<ApiResponse<null>>(`/products/notifications/${activityId}`, params);
  },

  // Delete a notification
  delete: async (activityId: string): Promise<void> => {
    await api.delete<ApiResponse<null>>(`/products/notifications/${activityId}`);
  },
};
