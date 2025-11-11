import client, { ApiResponse } from './client';

export interface SystemInfo {
  hostname: string;
  platform: string;
  os: string;
  architecture: string;
  cpuCores: number;
  uptime: number;
  bootTime: number;
}

export interface SystemMetrics {
  cpu: {
    usagePercent: number;
    perCore?: number[];
  };
  memory: {
    total: number;
    available: number;
    used: number;
    usedPercent: number;
  };
  disk: Array<{
    device: string;
    mountpoint: string;
    fstype: string;
    total: number;
    free: number;
    used: number;
    usedPercent: number;
  }>;
  network: {
    bytesSent: number;
    bytesRecv: number;
    packetsSent: number;
    packetsRecv: number;
  };
  timestamp: number;
}

export const systemApi = {
  getInfo: async () => {
    const response = await client.get<ApiResponse<SystemInfo>>('/system/info');
    return response.data;
  },

  getMetrics: async () => {
    const response = await client.get<ApiResponse<SystemMetrics>>('/system/metrics');
    return response.data;
  },
};
