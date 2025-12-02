import client, { ApiResponse } from './client';

// Types
export interface UPSConfig {
  id: number;
  created_at: string;
  updated_at: string;
  enabled: boolean;
  ups_name: string;
  ups_host: string;
  ups_port: number;
  ups_username?: string;
  ups_password?: string;
  poll_interval: number;
  low_battery_shutdown: boolean;
  low_battery_threshold: number;
  shutdown_delay: number;
  shutdown_command: string;
  notify_on_power_loss: boolean;
  notify_on_battery_low: boolean;
  notify_on_power_restored: boolean;
}

export interface UPSStatus {
  online: boolean;
  battery_charge: number;
  runtime: number; // seconds
  load_percent: number;
  input_voltage: number;
  output_voltage: number;
  temperature: number;
  status: string;
  model: string;
  manufacturer: string;
  serial: string;
  last_update: string;
}

export interface UPSEvent {
  id: number;
  created_at: string;
  event_type: string;
  description: string;
  battery_level?: number;
  runtime?: number;
  load_percent?: number;
  voltage?: number;
  severity: string;
}

// API
export const upsApi = {
  // Configuration
  async getConfig(): Promise<ApiResponse<UPSConfig>> {
    const response = await client.get('/ups/config');
    return response.data;
  },

  async updateConfig(config: Partial<UPSConfig>): Promise<ApiResponse<UPSConfig>> {
    const response = await client.put('/ups/config', config);
    return response.data;
  },

  async testConnection(config: Partial<UPSConfig>): Promise<ApiResponse<UPSStatus>> {
    const response = await client.post('/ups/test', config);
    return response.data;
  },

  // Status
  async getStatus(): Promise<ApiResponse<UPSStatus>> {
    const response = await client.get('/ups/status');
    return response.data;
  },

  // Events
  async getEvents(limit: number = 100, offset: number = 0): Promise<ApiResponse<UPSEvent[]>> {
    const response = await client.get('/ups/events', {
      params: { limit, offset },
    });
    return response.data;
  },

  // Monitoring
  async startMonitoring(): Promise<ApiResponse<any>> {
    const response = await client.post('/ups/monitoring/start');
    return response.data;
  },

  async stopMonitoring(): Promise<ApiResponse<any>> {
    const response = await client.post('/ups/monitoring/stop');
    return response.data;
  },
};
