import client, { ApiResponse } from './client';

// Types
export interface DockerContainer {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  created: string;
  ports?: Array<{
    privatePort: number;
    publicPort?: number;
    type: string;
  }>;
  labels?: Record<string, string>;
  networks?: Record<string, any>;
}

export interface DockerImage {
  id: string;
  repoTags: string[];
  repoDigests: string[];
  created: number;
  size: number;
  virtualSize: number;
  sharedSize: number;
  labels?: Record<string, string>;
  containers: number;
}

export interface DockerVolume {
  name: string;
  driver: string;
  mountpoint: string;
  createdAt: string;
  labels?: Record<string, string>;
  scope: string;
  options?: Record<string, string>;
}

export interface DockerNetwork {
  id: string;
  name: string;
  driver: string;
  scope: string;
  internal: boolean;
  enableIPv6: boolean;
  ipam?: {
    driver: string;
    config?: Array<{
      subnet?: string;
      gateway?: string;
    }>;
  };
  containers?: Record<string, any>;
  options?: Record<string, string>;
  labels?: Record<string, string>;
}

export interface PullImageRequest {
  image: string;
}

export interface ComposeService {
  name: string;
  image: string;
  status: string;
  containers: string[];
}

export interface ComposeStack {
  name: string;
  path: string;
  status: string;
  services: ComposeService[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateStackRequest {
  name: string;
  compose: string;
}

export interface UpdateStackRequest {
  compose: string;
}

export interface CreateContainerRequest {
  name: string;
  image: string;
  command?: string[];
  env?: string[];
  ports?: Array<{
    container: number;
    host?: number;
    protocol?: string;
  }>;
  volumes?: Array<{
    host: string;
    container: string;
    mode?: string;
  }>;
  restart?: string;
  network?: string;
  labels?: Record<string, string>;
}

// API
export const dockerApi = {
  // Containers
  async listContainers(all: boolean = false): Promise<ApiResponse<DockerContainer[]>> {
    const response = await client.get('/docker/containers', {
      params: { all },
    });
    return response.data;
  },

  async startContainer(id: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/containers/${id}/start`);
    return response.data;
  },

  async stopContainer(id: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/containers/${id}/stop`);
    return response.data;
  },

  async restartContainer(id: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/containers/${id}/restart`);
    return response.data;
  },

  async removeContainer(id: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/docker/containers/${id}`);
    return response.data;
  },

  async getContainerLogs(id: string): Promise<ApiResponse<string>> {
    const response = await client.get(`/docker/containers/${id}/logs`);
    return response.data;
  },

  async getContainerTop(id: string): Promise<ApiResponse<any>> {
    const response = await client.get(`/docker/containers/${id}/top`);
    return response.data;
  },

  async execContainer(id: string, command: string[]): Promise<ApiResponse<{ output: string }>> {
    const response = await client.post(`/docker/containers/${id}/exec`, { command });
    return response.data;
  },

  async updateContainerResources(
    id: string,
    resources: {
      memory?: number;
      memorySwap?: number;
      cpuShares?: number;
      cpuQuota?: number;
      cpuPeriod?: number;
    }
  ): Promise<ApiResponse<any>> {
    const response = await client.put(`/docker/containers/${id}/resources`, resources);
    return response.data;
  },

  // Images
  async listImages(): Promise<ApiResponse<DockerImage[]>> {
    const response = await client.get('/docker/images');
    return response.data;
  },

  async pullImage(image: string): Promise<ApiResponse<string>> {
    const response = await client.post('/docker/images/pull', { image });
    return response.data;
  },

  async removeImage(id: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/docker/images/${id}`);
    return response.data;
  },

  async buildImage(
    dockerfile: string,
    tags: string[],
    buildArgs?: Record<string, string>,
    labels?: Record<string, string>
  ): Promise<ApiResponse<{ output: string }>> {
    const response = await client.post('/docker/images/build', {
      dockerfile,
      tags,
      buildArgs,
      labels,
    });
    return response.data;
  },

  async tagImage(id: string, repo: string, tag: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/images/${id}/tag`, { repo, tag });
    return response.data;
  },

  async pushImage(id: string, registryAuth?: string): Promise<ApiResponse<{ output: string }>> {
    const response = await client.post(`/docker/images/${id}/push`, { registryAuth });
    return response.data;
  },

  // Volumes
  async listVolumes(): Promise<ApiResponse<DockerVolume[]>> {
    const response = await client.get('/docker/volumes');
    return response.data;
  },

  async removeVolume(id: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/docker/volumes/${id}`);
    return response.data;
  },

  // Networks
  async listNetworks(): Promise<ApiResponse<DockerNetwork[]>> {
    const response = await client.get('/docker/networks');
    return response.data;
  },

  // Stacks
  async listStacks(): Promise<ApiResponse<ComposeStack[]>> {
    const response = await client.get('/docker/stacks');
    return response.data;
  },

  async getStack(name: string): Promise<ApiResponse<ComposeStack>> {
    const response = await client.get(`/docker/stacks/${name}`);
    return response.data;
  },

  async createStack(request: CreateStackRequest): Promise<ApiResponse<any>> {
    const response = await client.post('/docker/stacks', request);
    return response.data;
  },

  async updateStack(name: string, request: UpdateStackRequest): Promise<ApiResponse<any>> {
    const response = await client.put(`/docker/stacks/${name}`, request);
    return response.data;
  },

  async deleteStack(name: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/docker/stacks/${name}`);
    return response.data;
  },

  async deployStack(name: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/stacks/${name}/deploy`);
    return response.data;
  },

  async stopStack(name: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/stacks/${name}/stop`);
    return response.data;
  },

  async restartStack(name: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/stacks/${name}/restart`);
    return response.data;
  },

  async removeStack(name: string, removeVolumes: boolean = false): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/stacks/${name}/remove`, null, {
      params: { volumes: removeVolumes },
    });
    return response.data;
  },

  async getStackLogs(name: string): Promise<ApiResponse<string>> {
    const response = await client.get(`/docker/stacks/${name}/logs`);
    return response.data;
  },

  async getStackCompose(name: string): Promise<ApiResponse<string>> {
    const response = await client.get(`/docker/stacks/${name}/compose`);
    return response.data;
  },

  async createContainer(request: CreateContainerRequest): Promise<ApiResponse<any>> {
    const response = await client.post('/docker/containers', request);
    return response.data;
  },

  // Template Management
  async listTemplates(category?: string): Promise<ApiResponse<any[]>> {
    const params = category ? { category } : undefined;
    const response = await client.get('/docker/templates', { params });
    return response.data;
  },

  async getTemplate(id: string): Promise<ApiResponse<any>> {
    const response = await client.get(`/docker/templates/${id}`);
    return response.data;
  },

  async getTemplateCategories(): Promise<ApiResponse<string[]>> {
    const response = await client.get('/docker/templates/categories');
    return response.data;
  },

  async deployTemplate(
    templateId: string,
    stackName: string,
    variables: Record<string, string>
  ): Promise<ApiResponse<any>> {
    const response = await client.post(`/docker/templates/${templateId}/deploy`, {
      stack_name: stackName,
      variables,
    });
    return response.data;
  },
};
