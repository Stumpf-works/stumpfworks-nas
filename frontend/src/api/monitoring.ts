import client, { ApiResponse } from './client';

export interface MonitoringConfig {
  prometheus_enabled: boolean;
  grafana_url: string;
  datadog_enabled: boolean;
  datadog_api_key?: string;
  datadog_api_key_set?: boolean;
}

export interface SystemMetrics {
  timestamp: string;
  cpu_usage_percent: number;
  memory_usage_percent: number;
  memory_total_mb: number;
  memory_used_mb: number;
  memory_available_mb: number;
  disk_usage_percent: number;
  disk_total_gb: number;
  disk_used_gb: number;
  disk_available_gb: number;
  network_rx_bytes: number;
  network_tx_bytes: number;
  load_average_1m: number;
  load_average_5m: number;
  load_average_15m: number;
  uptime_seconds: number;
}

export interface HealthScore {
  timestamp: string;
  score: number;
  status: 'healthy' | 'warning' | 'critical';
  details: {
    cpu_health: number;
    memory_health: number;
    disk_health: number;
    network_health: number;
  };
}

export interface MetricsTrend {
  metric: string;
  trend: 'increasing' | 'decreasing' | 'stable';
  change_percent: number;
  current_value: number;
  previous_value: number;
}

export const monitoringApi = {
  // Configuration
  getConfig: async () => {
    const response = await client.get<ApiResponse<MonitoringConfig>>('/monitoring/config');
    return response.data;
  },

  updateConfig: async (config: MonitoringConfig) => {
    const response = await client.put<ApiResponse<{ message: string }>>('/monitoring/config', config);
    return response.data;
  },

  // Metrics
  getLatestMetrics: async () => {
    const response = await client.get<ApiResponse<SystemMetrics>>('/metrics/latest');
    return response.data;
  },

  getMetricsHistory: async (params?: { start?: string; end?: string; limit?: number }) => {
    const queryParams = new URLSearchParams();
    if (params?.start) queryParams.append('start', params.start);
    if (params?.end) queryParams.append('end', params.end);
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const response = await client.get<ApiResponse<{ metrics: SystemMetrics[]; start: string; end: string; count: number }>>(
      `/metrics/history?${queryParams.toString()}`
    );
    return response.data;
  },

  getTrends: async (duration?: string) => {
    const queryParams = duration ? `?duration=${duration}` : '';
    const response = await client.get<ApiResponse<{ trends: Record<string, MetricsTrend>; duration: string }>>(
      `/metrics/trends${queryParams}`
    );
    return response.data;
  },

  // Health
  getLatestHealthScore: async () => {
    const response = await client.get<ApiResponse<HealthScore>>('/health/score');
    return response.data;
  },

  getHealthScores: async (params?: { start?: string; end?: string; limit?: number }) => {
    const queryParams = new URLSearchParams();
    if (params?.start) queryParams.append('start', params.start);
    if (params?.end) queryParams.append('end', params.end);
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const response = await client.get<ApiResponse<{ scores: HealthScore[]; start: string; end: string; count: number }>>(
      `/health/scores?${queryParams.toString()}`
    );
    return response.data;
  },
};
