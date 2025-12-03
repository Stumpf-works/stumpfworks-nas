import client, { ApiResponse } from './client';

// Types
export interface TimeMachineDevice {
  id: number;
  created_at: string;
  updated_at: string;
  device_name: string;
  mac_address?: string;
  model_id?: string;
  share_path: string;
  quota_gb: number;
  used_gb: number;
  enabled: boolean;
  last_backup?: string;
  last_seen?: string;
  username?: string;
}

export interface TimeMachineConfig {
  id?: number;
  created_at?: string;
  updated_at?: string;
  enabled: boolean;
  share_name: string;
  base_path: string;
  default_quota_gb: number;
  auto_discovery: boolean;
  use_afp: boolean;
  use_smb: boolean;
  smb_version: string;
}

export interface CreateDeviceRequest {
  device_name: string;
  mac_address?: string;
  model_id?: string;
  share_path?: string;
  quota_gb?: number;
  username?: string;
  password?: string;
  enabled?: boolean;
}

export interface UpdateDeviceRequest {
  device_name?: string;
  mac_address?: string;
  model_id?: string;
  share_path?: string;
  quota_gb?: number;
  username?: string;
  password?: string;
  enabled?: boolean;
}

class TimeMachineAPI {
  // Configuration
  async getConfig(): Promise<ApiResponse<TimeMachineConfig>> {
    const response = await client.get('/timemachine/config');
    return response.data;
  }

  async updateConfig(config: Partial<TimeMachineConfig>): Promise<ApiResponse<TimeMachineConfig>> {
    const response = await client.put('/timemachine/config', config);
    return response.data;
  }

  async enable(): Promise<ApiResponse<{ success: boolean; message: string }>> {
    const response = await client.post('/timemachine/enable');
    return response.data;
  }

  async disable(): Promise<ApiResponse<{ success: boolean; message: string }>> {
    const response = await client.post('/timemachine/disable');
    return response.data;
  }

  // Devices
  async listDevices(): Promise<ApiResponse<TimeMachineDevice[]>> {
    const response = await client.get('/timemachine/devices');
    return response.data;
  }

  async getDevice(id: number): Promise<ApiResponse<TimeMachineDevice>> {
    const response = await client.get(`/timemachine/devices/${id}`);
    return response.data;
  }

  async createDevice(device: CreateDeviceRequest): Promise<ApiResponse<TimeMachineDevice>> {
    const response = await client.post('/timemachine/devices', device);
    return response.data;
  }

  async updateDevice(id: number, device: UpdateDeviceRequest): Promise<ApiResponse<TimeMachineDevice>> {
    const response = await client.put(`/timemachine/devices/${id}`, device);
    return response.data;
  }

  async deleteDevice(id: number): Promise<ApiResponse<{ success: boolean; message: string }>> {
    const response = await client.delete(`/timemachine/devices/${id}`);
    return response.data;
  }

  // Usage monitoring
  async updateDeviceUsage(id: number): Promise<ApiResponse<TimeMachineDevice>> {
    const response = await client.post(`/timemachine/devices/${id}/update-usage`);
    return response.data;
  }

  async updateAllDeviceUsages(): Promise<ApiResponse<TimeMachineDevice[]>> {
    const response = await client.post('/timemachine/devices/update-all-usage');
    return response.data;
  }
}

export const timeMachineApi = new TimeMachineAPI();
