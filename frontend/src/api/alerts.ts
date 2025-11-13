import client, { ApiResponse } from './client';

export interface AlertConfig {
  id?: number;
  createdAt?: string;
  updatedAt?: string;

  // Email settings
  enabled: boolean;
  smtpHost: string;
  smtpPort: number;
  smtpUsername: string;
  smtpPassword?: string;
  smtpFromEmail: string;
  smtpFromName: string;
  smtpUseTLS: boolean;
  alertRecipient: string;

  // Webhook settings
  webhookEnabled: boolean;
  webhookType: 'discord' | 'slack' | 'custom';
  webhookURL: string;
  webhookUsername: string;
  webhookAvatarURL: string;

  // Alert triggers
  onFailedLogin: boolean;
  onIPBlock: boolean;
  onCriticalEvent: boolean;
  failedLoginThreshold: number;

  // Rate limiting
  rateLimitMinutes: number;
}

export interface AlertLog {
  id: number;
  createdAt: string;
  alertType: string;
  channel: 'email' | 'webhook';
  subject: string;
  body: string;
  recipient: string;
  status: 'sent' | 'failed';
  error?: string;
}

export const alertsApi = {
  // Get alert configuration
  getConfig: async () => {
    const response = await client.get<ApiResponse<AlertConfig>>('/alerts/config');
    return response.data;
  },

  // Update alert configuration
  updateConfig: async (config: AlertConfig) => {
    const response = await client.put<ApiResponse<AlertConfig>>('/alerts/config', config);
    return response.data;
  },

  // Test email configuration
  testEmail: async (config: AlertConfig) => {
    const response = await client.post<ApiResponse<{ message: string }>>('/alerts/test/email', config);
    return response.data;
  },

  // Test webhook configuration
  testWebhook: async (config: AlertConfig) => {
    const response = await client.post<ApiResponse<{ message: string }>>('/alerts/test/webhook', config);
    return response.data;
  },

  // Get alert logs
  getLogs: async (limit = 50) => {
    const response = await client.get<ApiResponse<AlertLog[]>>(`/alerts/logs?limit=${limit}`);
    return response.data;
  },
};
