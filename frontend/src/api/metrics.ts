import client from './client';

export interface SystemMetric {
  id: number;
  timestamp: string;

  // CPU metrics
  cpuUsage: number;
  cpuLoadAvg1: number;
  cpuLoadAvg5: number;
  cpuLoadAvg15: number;
  cpuTemperature: number;

  // Memory metrics
  memoryUsedBytes: number;
  memoryTotalBytes: number;
  memoryUsage: number;
  swapUsedBytes: number;
  swapTotalBytes: number;
  swapUsage: number;

  // Disk metrics
  diskUsedBytes: number;
  diskTotalBytes: number;
  diskUsage: number;
  diskReadBytesPerSec: number;
  diskWriteBytesPerSec: number;
  diskIOPS: number;

  // Network metrics
  networkRxBytesPerSec: number;
  networkTxBytesPerSec: number;
  networkRxPacketsPerSec: number;
  networkTxPacketsPerSec: number;

  // Process metrics
  processCount: number;
  threadCount: number;

  createdAt: string;
}

export interface HealthScore {
  id: number;
  timestamp: string;
  score: number;
  cpuScore: number;
  memoryScore: number;
  diskScore: number;
  networkScore: number;
  issues?: string;
  createdAt: string;
}

export interface MetricsTrend {
  metricName: string;
  currentValue: number;
  previousValue: number;
  change: number;
  changePercent: number;
  direction: 'up' | 'down' | 'stable';
  timestamp: string;
}

export interface MetricsHistoryResponse {
  metrics: SystemMetric[];
  start: string;
  end: string;
  count: number;
}

export interface HealthScoresResponse {
  scores: HealthScore[];
  start: string;
  end: string;
  count: number;
}

export interface TrendsResponse {
  trends: MetricsTrend[];
  duration: string;
}

export const metricsApi = {
  /**
   * Get historical metrics
   */
  getHistory: async (
    start?: string,
    end?: string,
    limit?: number
  ): Promise<MetricsHistoryResponse> => {
    const params: any = {};
    if (start) params.start = start;
    if (end) params.end = end;
    if (limit) params.limit = limit;

    const response = await client.get('/metrics/history', { params });
    return response.data.data;
  },

  /**
   * Get latest metric
   */
  getLatest: async (): Promise<SystemMetric> => {
    const response = await client.get('/metrics/latest');
    return response.data.data;
  },

  /**
   * Get trends
   */
  getTrends: async (duration?: string): Promise<TrendsResponse> => {
    const params: any = {};
    if (duration) params.duration = duration;

    const response = await client.get('/metrics/trends', { params });
    return response.data.data;
  },

  /**
   * Get health scores
   */
  getHealthScores: async (
    start?: string,
    end?: string,
    limit?: number
  ): Promise<HealthScoresResponse> => {
    const params: any = {};
    if (start) params.start = start;
    if (end) params.end = end;
    if (limit) params.limit = limit;

    const response = await client.get('/health/scores', { params });
    return response.data.data;
  },

  /**
   * Get latest health score
   */
  getLatestHealthScore: async (): Promise<HealthScore> => {
    const response = await client.get('/health/score');
    return response.data.data;
  },
};
