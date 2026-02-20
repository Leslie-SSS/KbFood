// Product entity
export interface Product {
  id: number;
  activityId: string;
  platform: string;
  region: string;
  title: string;
  shopName: string;
  originalPrice: number;
  currentPrice: number;
  salesStatus: number;
  salesStatusText: string;
  activityCreateTime: string;
  createTime: string;
  updateTime: string;
  discount?: number;
  dropRate?: number;
  hasNotification?: boolean;
  targetPrice?: number | null;
}

// Product filter parameters
export interface ProductFilters {
  keyword: string;
  platform: string;
  region: string;
  salesStatus: string;
  monitorStatus: string;
  recentSevenDays: boolean;
}

// Platform options for filter dropdown
export const PLATFORMS = [
  { value: '', label: '所有平台' },
  { value: '探探糖', label: '探探糖' },
  { value: 'DT', label: 'DT' },
  { value: '小蚕', label: '小蚕' },
] as const;

// Region options for filter dropdown
export const REGIONS = [
  { value: '', label: '所有地区' },
  { value: '广州', label: '广州' },
  { value: '佛山', label: '佛山' },
] as const;

// Sales status options for filter dropdown
export const SALES_STATUS = [
  { value: '', label: '销售状态' },
  { value: '1', label: '在售' },
  { value: '0', label: '已售罄' },
] as const;

// Monitor status options for filter dropdown
export const MONITOR_STATUS = [
  { value: '', label: '监控状态' },
  { value: '1', label: '已监控' },
  { value: '0', label: '未监控' },
] as const;
