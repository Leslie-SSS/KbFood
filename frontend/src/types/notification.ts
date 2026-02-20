// Notification configuration
export interface NotificationConfig {
  activityId: string;
  targetPrice: number;
  lastNotifyTime?: string | null;
}

// Parameters for creating a notification
export interface CreateNotificationParams {
  activityId: string;
  targetPrice: number;
}

// Parameters for updating a notification
export interface UpdateNotificationParams {
  targetPrice: number;
}
