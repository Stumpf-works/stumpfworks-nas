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

// Pacemaker/Corosync Types
export interface ClusterNode {
  name: string;
  id: number;
  ip: string;
  online: boolean;
}

export interface ClusterResource {
  id: string;
  type: string;
  agent: string;
  node: string;
  active: boolean;
  managed: boolean;
  failed: boolean;
}

export interface ClusterStatus {
  name: string;
  nodes: ClusterNode[];
  resources: ClusterResource[];
  quorum: boolean;
  maintenance_mode: boolean;
  stonith_enabled: boolean;
  symmetric_cluster: boolean;
}

export interface OpConfig {
  name: string;
  interval?: string;
  timeout?: string;
}

export interface ResourceConfig {
  id: string;
  type?: string;
  agent: string;
  params?: Record<string, string>;
  op?: OpConfig[];
}

// Keepalived (VIP) Types
export interface VIPConfig {
  id?: string;
  virtual_ip: string;
  interface: string;
  router_id?: number;
  priority?: number;
  state?: string;
  auth_pass?: string;
  virtual_routes?: string[];
  track_scripts?: string[];
}

export interface VIPStatus {
  id: string;
  virtual_ip: string;
  interface: string;
  state: string;
  is_master: boolean;
  priority: number;
  is_active: boolean;
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

  // Pacemaker/Corosync Cluster Management
  getClusterStatus: async (): Promise<ApiResponse<ClusterStatus>> => {
    const response = await fetch('/api/v1/ha/cluster/status');
    return response.json();
  },

  createClusterResource: async (config: ResourceConfig): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/ha/cluster/resources', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config),
    });
    return response.json();
  },

  deleteClusterResource: async (resourceId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/resources/${encodeURIComponent(resourceId)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  enableClusterResource: async (resourceId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/resources/${encodeURIComponent(resourceId)}/enable`, {
      method: 'POST',
    });
    return response.json();
  },

  disableClusterResource: async (resourceId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/resources/${encodeURIComponent(resourceId)}/disable`, {
      method: 'POST',
    });
    return response.json();
  },

  moveClusterResource: async (resourceId: string, targetNode: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/resources/${encodeURIComponent(resourceId)}/move`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_node: targetNode }),
    });
    return response.json();
  },

  clearClusterResource: async (resourceId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/resources/${encodeURIComponent(resourceId)}/clear`, {
      method: 'POST',
    });
    return response.json();
  },

  setMaintenanceMode: async (enabled: boolean): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/ha/cluster/maintenance', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled }),
    });
    return response.json();
  },

  standbyNode: async (nodeName: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/nodes/${encodeURIComponent(nodeName)}/standby`, {
      method: 'POST',
    });
    return response.json();
  },

  unstandbyNode: async (nodeName: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/cluster/nodes/${encodeURIComponent(nodeName)}/unstandby`, {
      method: 'POST',
    });
    return response.json();
  },

  // Keepalived Virtual IP Management
  listVIPs: async (): Promise<ApiResponse<VIPStatus[]>> => {
    const response = await fetch('/api/v1/ha/vip/');
    return response.json();
  },

  getVIPStatus: async (vipId: string): Promise<ApiResponse<VIPStatus>> => {
    const response = await fetch(`/api/v1/ha/vip/${encodeURIComponent(vipId)}`);
    return response.json();
  },

  createVIP: async (config: VIPConfig): Promise<ApiResponse<void>> => {
    const response = await fetch('/api/v1/ha/vip/', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config),
    });
    return response.json();
  },

  deleteVIP: async (vipId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/vip/${encodeURIComponent(vipId)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  promoteVIPToMaster: async (vipId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/vip/${encodeURIComponent(vipId)}/promote`, {
      method: 'POST',
    });
    return response.json();
  },

  demoteVIPToBackup: async (vipId: string): Promise<ApiResponse<void>> => {
    const response = await fetch(`/api/v1/ha/vip/${encodeURIComponent(vipId)}/demote`, {
      method: 'POST',
    });
    return response.json();
  },
};
