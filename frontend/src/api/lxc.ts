import { ApiResponse } from './client';

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
    const response = await fetch('/api/v1/lxc/containers');
    return response.json();
  },

  // Get container details
  getContainer: async (name: string): Promise<ApiResponse<Container>> => {
    const response = await fetch(`/api/v1/lxc/containers/${encodeURIComponent(name)}`);
    return response.json();
  },

  // Create a new container
  createContainer: async (data: ContainerCreateRequest): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch('/api/v1/lxc/containers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Start a container
  startContainer: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch(`/api/v1/lxc/containers/${encodeURIComponent(name)}/start`, {
      method: 'POST',
    });
    return response.json();
  },

  // Stop a container
  stopContainer: async (name: string, force: boolean = false): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch(`/api/v1/lxc/containers/${encodeURIComponent(name)}/stop?force=${force}`, {
      method: 'POST',
    });
    return response.json();
  },

  // Delete a container
  deleteContainer: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch(`/api/v1/lxc/containers/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  // List available templates
  listTemplates: async (): Promise<ApiResponse<LXCTemplate[]>> => {
    const response = await fetch('/api/v1/lxc/templates');
    return response.json();
  },
};
