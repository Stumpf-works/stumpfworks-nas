import { ApiResponse } from './client';

// DRBD Types
export interface DRBDResource {
  name: string;
  device: string;        // e.g., /dev/drbd0
  disk: string;          // e.g., /dev/sda1
  meta_disk: string;     // e.g., internal or /dev/sdb1
  local_address: string; // e.g., 192.168.1.10:7788
  peer_address: string;  // e.g., 192.168.1.11:7788
  protocol: string;      // A, B, or C
}

export interface DRBDStatus {
  name: string;
  device: string;
  connection_state: string; // Connected, Disconnected, StandAlone
  role: string;             // Primary, Secondary, Unknown
  disk_state: string;       // UpToDate, Inconsistent, DUnknown
  peer_role: string;        // Primary, Secondary, Unknown
  peer_disk_state: string;  // UpToDate, Inconsistent, DUnknown
  sync_progress: number;    // 0-100 percentage
  resyncing: boolean;
}

export interface CreateDRBDResourceRequest {
  name: string;
  device: string;
  disk: string;
  meta_disk?: string;
  local_address: string;
  peer_address: string;
  protocol?: string;
}

export const haApi = {
  // DRBD Resource Management
  listDRBDResources: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ha/drbd/resources');
    return response.json();
  },

  getDRBDResourceStatus: async (name: string): Promise<ApiResponse<DRBDStatus>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}`);
    return response.json();
  },

  createDRBDResource: async (request: CreateDRBDResourceRequest): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/ha/drbd/resources', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
    return response.json();
  },

  deleteDRBDResource: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  // DRBD Resource Operations
  promoteDRBDResource: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/promote`, {
      method: 'POST',
    });
    return response.json();
  },

  demoteDRBDResource: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/demote`, {
      method: 'POST',
    });
    return response.json();
  },

  forcePrimaryDRBDResource: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/force-primary`, {
      method: 'POST',
    });
    return response.json();
  },

  disconnectDRBDResource: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/disconnect`, {
      method: 'POST',
    });
    return response.json();
  },

  connectDRBDResource: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/connect`, {
      method: 'POST',
    });
    return response.json();
  },

  startDRBDSync: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/sync`, {
      method: 'POST',
    });
    return response.json();
  },

  verifyDRBDData: async (name: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/drbd/resources/${encodeURIComponent(name)}/verify`, {
      method: 'POST',
    });
    return response.json();
  },
};
