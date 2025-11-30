import client, { ApiResponse } from './client';

// LXC Container Types
export interface Container {
  name: string;
  state: string;
  ipv4: string;
  ipv6: string;
  autostart: boolean;
  pid: number;
  memory: string;
  memory_limit: string;
  cpu_usage: string;
}

export interface ContainerCreateRequest {
  name: string;
  template: string;
  release?: string;
  arch?: string;
  storage?: string;
  network?: string;
  autostart?: boolean;
}

export interface LXCTemplate {
  name: string;
  description: string;
}

export const lxcApi = {
  // List all containers
  listContainers: async (): Promise<ApiResponse<Container[]>> => {
    const response = await client.get<ApiResponse<Container[]>>('/lxc/containers');
    return response.data;
  },

  // Get container details
  getContainer: async (name: string): Promise<ApiResponse<Container>> => {
    const response = await client.get<ApiResponse<Container>>(`/lxc/containers/${encodeURIComponent(name)}`);
    return response.data;
  },

  // Create a new container
  createContainer: async (data: ContainerCreateRequest): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>('/lxc/containers', data);
    return response.data;
  },

  // Start a container
  startContainer: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>(`/lxc/containers/${encodeURIComponent(name)}/start`);
    return response.data;
  },

  // Stop a container
  stopContainer: async (name: string, force: boolean = false): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>(`/lxc/containers/${encodeURIComponent(name)}/stop?force=${force}`);
    return response.data;
  },

  // Delete a container
  deleteContainer: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await client.delete<ApiResponse<{ message: string; name: string }>>(`/lxc/containers/${encodeURIComponent(name)}`);
    return response.data;
  },

  // List available templates
  listTemplates: async (): Promise<ApiResponse<LXCTemplate[]>> => {
    const response = await client.get<ApiResponse<LXCTemplate[]>>('/lxc/templates');
    return response.data;
  },
};
