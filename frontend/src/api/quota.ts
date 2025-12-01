import { ApiResponse } from './client';

export interface QuotaInfo {
  name: string;
  type: 'user' | 'group';
  filesystem: string;
  blocksUsed: number;
  blocksSoft: number;
  blocksHard: number;
  blocksGrace: string;
  filesUsed: number;
  filesSoft: number;
  filesHard: number;
  filesGrace: string;
}

export interface SetQuotaRequest {
  name: string;
  filesystem: string;
  blocksSoft?: number;
  blocksHard?: number;
  filesSoft?: number;
  filesHard?: number;
}

export interface RemoveQuotaRequest {
  name: string;
  filesystem: string;
}

export interface FilesystemQuotaStatus {
  filesystem: string;
  quotasEnabled: boolean;
  userQuotaEnabled: boolean;
  groupQuotaEnabled: boolean;
}

export const quotaApi = {
  // User quota operations
  getUserQuota: async (username: string, filesystem: string): Promise<ApiResponse<QuotaInfo>> => {
    const response = await fetch(`/api/v1/quotas/user?name=${encodeURIComponent(username)}&filesystem=${encodeURIComponent(filesystem)}`);
    return response.json();
  },

  setUserQuota: async (request: SetQuotaRequest): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/quotas/user', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
    return response.json();
  },

  removeUserQuota: async (request: RemoveQuotaRequest): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/quotas/user', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
    return response.json();
  },

  listUserQuotas: async (filesystem: string): Promise<ApiResponse<QuotaInfo[]>> => {
    const response = await fetch(`/api/v1/quotas/users?filesystem=${encodeURIComponent(filesystem)}`);
    return response.json();
  },

  // Group quota operations
  getGroupQuota: async (groupname: string, filesystem: string): Promise<ApiResponse<QuotaInfo>> => {
    const response = await fetch(`/api/v1/quotas/group?name=${encodeURIComponent(groupname)}&filesystem=${encodeURIComponent(filesystem)}`);
    return response.json();
  },

  setGroupQuota: async (request: SetQuotaRequest): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/quotas/group', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
    return response.json();
  },

  removeGroupQuota: async (request: RemoveQuotaRequest): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/quotas/group', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
    return response.json();
  },

  listGroupQuotas: async (filesystem: string): Promise<ApiResponse<QuotaInfo[]>> => {
    const response = await fetch(`/api/v1/quotas/groups?filesystem=${encodeURIComponent(filesystem)}`);
    return response.json();
  },

  // Filesystem quota status
  getQuotaStatus: async (filesystem: string): Promise<ApiResponse<FilesystemQuotaStatus>> => {
    const response = await fetch(`/api/v1/quotas/status?filesystem=${encodeURIComponent(filesystem)}`);
    return response.json();
  },
};
