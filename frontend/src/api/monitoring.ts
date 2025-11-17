import client, { ApiResponse } from './client';

export interface MonitoringConfig {
  prometheus_enabled: boolean;
  grafana_url: string;
  datadog_enabled: boolean;
  datadog_api_key?: string;
  datadog_api_key_set?: boolean;
}

export const monitoringApi = {
  // Get monitoring configuration
  getConfig: async () => {
    const response = await client.get<ApiResponse<MonitoringConfig>>('/monitoring/config');
    return response.data;
  },

  // Update monitoring configuration
  updateConfig: async (config: MonitoringConfig) => {
    const response = await client.put<ApiResponse<{ message: string }>>('/monitoring/config', config);
    return response.data;
  },
};
