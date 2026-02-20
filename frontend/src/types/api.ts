// API Response wrapper
export interface ApiResponse<T> {
  code: number;
  data?: T;
  message?: string;
  meta?: {
    total: number;
    page: number;
    limit: number;
  };
}
