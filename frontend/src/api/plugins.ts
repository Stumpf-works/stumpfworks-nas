import client, { ApiResponse } from './client';

// Types
export interface Plugin {
  id: string;
  name: string;
  version: string;
  author: string;
  description: string;
  icon?: string;
  enabled: boolean;
  installed: boolean;
  installPath?: string;
  config?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface InstallPluginRequest {
  sourcePath: string;
}

export interface UpdatePluginConfigRequest {
  config: Record<string, any>;
}

// API
export const pluginApi = {
  // List all plugins
  async listPlugins(): Promise<ApiResponse<Plugin[]>> {
    const response = await client.get('/plugins');
    return response.data;
  },

  // Get a specific plugin
  async getPlugin(id: string): Promise<ApiResponse<Plugin>> {
    const response = await client.get(`/plugins/${id}`);
    return response.data;
  },

  // Install a plugin
  async installPlugin(request: InstallPluginRequest): Promise<ApiResponse<Plugin>> {
    const response = await client.post('/plugins/install', request);
    return response.data;
  },

  // Uninstall a plugin
  async uninstallPlugin(id: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/plugins/${id}`);
    return response.data;
  },

  // Enable a plugin
  async enablePlugin(id: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/plugins/${id}/enable`);
    return response.data;
  },

  // Disable a plugin
  async disablePlugin(id: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/plugins/${id}/disable`);
    return response.data;
  },

  // Update plugin configuration
  async updatePluginConfig(
    id: string,
    request: UpdatePluginConfigRequest
  ): Promise<ApiResponse<any>> {
    const response = await client.put(`/plugins/${id}/config`, request);
    return response.data;
  },
};
