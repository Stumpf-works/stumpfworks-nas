import client, { ApiResponse } from './client';

// Addon Types
export interface AddonManifest {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  version: string;
  author: string;
  system_packages: string[];
  services: string[];
  install_script?: string;
  uninstall_script?: string;
  app_component?: string;
  route_prefix?: string;
  minimum_memory: number;
  minimum_disk: number;
  architecture: string[];
}

export interface InstallationStatus {
  addon_id: string;
  installed: boolean;
  version: string;
  install_date: string;
  packages_ok: boolean;
  services_ok: boolean;
  error?: string;
}

export interface AddonWithStatus {
  manifest: AddonManifest;
  status: InstallationStatus;
}

export const addonsApi = {
  // List all available addons with their installation status
  listAddons: async (): Promise<ApiResponse<AddonWithStatus[]>> => {
    const response = await client.get<ApiResponse<AddonWithStatus[]>>('/addons/');
    return response.data;
  },

  // Get details of a specific addon
  getAddon: async (addonId: string): Promise<ApiResponse<AddonWithStatus>> => {
    const response = await client.get<ApiResponse<AddonWithStatus>>(`/addons/${encodeURIComponent(addonId)}`);
    return response.data;
  },

  // Get installation status of an addon
  getAddonStatus: async (addonId: string): Promise<ApiResponse<InstallationStatus>> => {
    const response = await client.get<ApiResponse<InstallationStatus>>(`/addons/${encodeURIComponent(addonId)}/status`);
    return response.data;
  },

  // Install an addon
  installAddon: async (addonId: string): Promise<ApiResponse<{ message: string; addon_id: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; addon_id: string }>>(`/addons/${encodeURIComponent(addonId)}/install`);
    return response.data;
  },

  // Uninstall an addon
  uninstallAddon: async (addonId: string): Promise<ApiResponse<{ message: string; addon_id: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; addon_id: string }>>(`/addons/${encodeURIComponent(addonId)}/uninstall`);
    return response.data;
  },
};
