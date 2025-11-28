import client, { ApiResponse } from './client';

export interface ADConfig {
  enabled: boolean;
  server: string;
  port: number;
  baseDN: string;
  bindUser: string;
  bindPassword?: string;
  userFilter: string;
  groupFilter: string;
  useTLS: boolean;
  skipVerify?: boolean;
}

export interface ADUser {
  username: string;
  displayName: string;
  email: string;
  groups: string[];
  distinguishedName: string;
  enabled: boolean;
}

export interface ADGroup {
  name: string;
  dn: string;
  members: string[];
}

export const adApi = {
  // Get AD configuration
  getConfig: async () => {
    const response = await client.get<ApiResponse<ADConfig>>('/ad/config');
    return response.data;
  },

  // Update AD configuration
  updateConfig: async (config: Partial<ADConfig>) => {
    const response = await client.put<ApiResponse<ADConfig>>('/ad/config', config);
    return response.data;
  },

  // Test AD connection
  testConnection: async () => {
    const response = await client.post<ApiResponse<{ success: boolean; message: string }>>('/ad/test');
    return response.data;
  },

  // Authenticate user
  authenticate: async (username: string, password: string) => {
    const response = await client.post<ApiResponse<ADUser>>('/ad/authenticate', { username, password });
    return response.data;
  },

  // List AD users
  listUsers: async () => {
    const response = await client.get<ApiResponse<ADUser[]>>('/ad/users');
    return response.data;
  },

  // Sync specific user
  syncUser: async (username: string) => {
    const response = await client.post<ApiResponse<ADUser>>('/ad/users/sync', { username });
    return response.data;
  },

  // Get user groups
  getUserGroups: async (username: string) => {
    const response = await client.get<ApiResponse<string[]>>(`/api/v1/ad/users/${username}/groups`);
    return response.data;
  },
};
