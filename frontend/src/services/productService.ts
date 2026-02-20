import { api } from './api';
import type { ApiResponse, Product, PriceTrend } from '@/types';

export const productService = {
  // Get products with filters
  getProducts: async (params: Record<string, string | boolean>): Promise<Product[]> => {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== '' && value !== false) {
        searchParams.append(key, String(value));
      }
    });

    const response = await api.get<ApiResponse<Product[]>>(
      `/products?${searchParams.toString()}`
    );
    return response.data.data || [];
  },

  // Get price trend for a product
  getPriceTrend: async (activityId: string): Promise<PriceTrend[]> => {
    const response = await api.get<ApiResponse<PriceTrend[]>>(
      `/products/${activityId}/trend`
    );
    return response.data.data || [];
  },

  // Block a product
  blockProduct: async (activityId: string): Promise<void> => {
    await api.post(`/products/${activityId}/block`);
  },

  // Unblock a product
  unblockProduct: async (activityId: string): Promise<void> => {
    await api.post(`/products/unblock/${activityId}`);
  },

  // Get all blocked products
  getBlockedProducts: async (): Promise<Product[]> => {
    const response = await api.get<ApiResponse<Product[]>>('/products/blocked');
    return response.data.data || [];
  },
};
