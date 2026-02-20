import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';

export interface SyncStatus {
  lastRunTime: string;
  status: 'success' | 'failed' | 'running' | 'pending';
  productCount: number;
  isHealthy: boolean;
  errorMessage?: string;
}

export interface SystemStatus {
  sync: SyncStatus;
  serverTime: string;
}

export function useSystemStatus() {
  return useQuery<SystemStatus>({
    queryKey: ['systemStatus'],
    queryFn: async () => {
      const response = await api.get('/status');
      return response.data.data;
    },
    refetchInterval: 60000, // Refresh every minute
    staleTime: 30000, // Consider data fresh for 30 seconds
  });
}
