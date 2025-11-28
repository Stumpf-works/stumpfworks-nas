// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
import client, { ApiResponse } from './client';

// ===== Types =====

// ZFS Types
export interface ZFSPool {
  name: string;
  health: string;
  size: number;
  allocated: number;
  free: number;
  capacity: number;
  dedup: number;
  fragmentation: number;
}

export interface ZFSDataset {
  name: string;
  type: string;
  used: number;
  available: number;
  referenced: number;
  mountpoint: string;
}

export interface ZFSSnapshot {
  name: string;
  creation: string;
  used: number;
  referenced: number;
}

export interface CreateZFSPoolRequest {
  name: string;
  raid_type: 'mirror' | 'raidz' | 'raidz2' | 'raidz3' | 'stripe';
  devices: string[];
  options?: Record<string, string>;
}

// RAID Types
export interface RAIDArray {
  name: string;
  level: string;
  devices: number;
  active: number;
  working: number;
  failed: number;
  spare: number;
  size: number;
  usedSize: number;
  state: string;
  resync: number;
}

export interface CreateRAIDArrayRequest {
  name: string;
  level: '0' | '1' | '5' | '6' | '10';
  devices: string[];
  spare?: string[];
}

// SMART Types
export interface SMARTInfo {
  device: string;
  smartStatus: string;
  temperature: number;
  powerOnHours: number;
  reallocatedSectors: number;
  pendingSectors: number;
  uncorrectableErrors: number;
  healthScore: number;
}

// Samba Types
export interface SambaShare {
  name: string;
  path: string;
  comment: string;
  validUsers: string[];
  validGroups: string[];
  readOnly: boolean;
  browseable: boolean;
  guestOK: boolean;
  recycleBin: boolean;
}

export interface CreateSambaShareRequest {
  name: string;
  path: string;
  comment?: string;
  valid_users?: string[];
  valid_groups?: string[];
  read_only?: boolean;
  browseable?: boolean;
  guest_ok?: boolean;
  recycle_bin?: boolean;
}

// NFS Types
export interface NFSExport {
  path: string;
  clients: string[];
  options: string[];
}

export interface CreateNFSExportRequest {
  path: string;
  clients: string[];
  options?: string[];
}

// Network Types
export interface CreateBondRequest {
  name: string;
  mode: string;
  interfaces: string[];
}

export interface CreateVLANRequest {
  parent: string;
  vlan_id: number;
}

// ===== API Client =====

export const syslibApi = {
  // System Library Health
  getHealth: async () => {
    const response = await client.get<ApiResponse<any>>('/syslib/health');
    return response.data;
  },

  // ZFS Operations
  zfs: {
    listPools: async () => {
      const response = await client.get<ApiResponse<ZFSPool[]>>('/syslib/zfs/pools');
      return response.data;
    },

    getPool: async (name: string) => {
      const response = await client.get<ApiResponse<ZFSPool>>(`/syslib/zfs/pools/${name}`);
      return response.data;
    },

    createPool: async (data: CreateZFSPoolRequest) => {
      const response = await client.post<ApiResponse<{ message: string; pool: string }>>('/syslib/zfs/pools', data);
      return response.data;
    },

    destroyPool: async (name: string, force: boolean = false) => {
      const response = await client.delete<ApiResponse<{ message: string }>>(`/syslib/zfs/pools/${name}?force=${force}`);
      return response.data;
    },

    scrubPool: async (name: string) => {
      const response = await client.post<ApiResponse<{ message: string }>>(`/syslib/zfs/pools/${name}/scrub`, {});
      return response.data;
    },

    listDatasets: async (poolName: string) => {
      const response = await client.get<ApiResponse<ZFSDataset[]>>(`/syslib/zfs/pools/${poolName}/datasets`);
      return response.data;
    },

    createSnapshot: async (dataset: string, snapshot: string) => {
      const response = await client.post<ApiResponse<{ message: string }>>('/syslib/zfs/snapshots', { dataset, snapshot });
      return response.data;
    },

    listSnapshots: async (dataset: string) => {
      const response = await client.get<ApiResponse<ZFSSnapshot[]>>(`/syslib/zfs/datasets/${dataset}/snapshots`);
      return response.data;
    },
  },

  // RAID Operations
  raid: {
    listArrays: async () => {
      const response = await client.get<ApiResponse<RAIDArray[]>>('/syslib/raid/arrays');
      return response.data;
    },

    getArray: async (name: string) => {
      const response = await client.get<ApiResponse<RAIDArray>>(`/syslib/raid/arrays/${name}`);
      return response.data;
    },

    createArray: async (data: CreateRAIDArrayRequest) => {
      const response = await client.post<ApiResponse<{ message: string; array: string }>>('/syslib/raid/arrays', data);
      return response.data;
    },
  },

  // SMART Operations
  smart: {
    getInfo: async (device: string) => {
      const response = await client.get<ApiResponse<SMARTInfo>>(`/syslib/smart/${device}`);
      return response.data;
    },

    runTest: async (device: string, type: 'short' | 'long' | 'conveyance' = 'short') => {
      const response = await client.post<ApiResponse<{ message: string; type: string }>>(`/syslib/smart/${device}/test?type=${type}`, {});
      return response.data;
    },
  },

  // Samba Operations
  samba: {
    getStatus: async () => {
      const response = await client.get<ApiResponse<{ active: boolean; enabled: boolean }>>('/syslib/samba/status');
      return response.data;
    },

    restart: async () => {
      const response = await client.post<ApiResponse<{ message: string }>>('/syslib/samba/restart', {});
      return response.data;
    },

    listShares: async () => {
      const response = await client.get<ApiResponse<SambaShare[]>>('/syslib/samba/shares');
      return response.data;
    },

    getShare: async (name: string) => {
      const response = await client.get<ApiResponse<SambaShare>>(`/syslib/samba/shares/${name}`);
      return response.data;
    },

    createShare: async (data: CreateSambaShareRequest) => {
      const response = await client.post<ApiResponse<{ message: string; share: string }>>('/syslib/samba/shares', data);
      return response.data;
    },

    deleteShare: async (name: string) => {
      const response = await client.delete<ApiResponse<{ message: string }>>(`/syslib/samba/shares/${name}`);
      return response.data;
    },
  },

  // NFS Operations
  nfs: {
    restart: async () => {
      const response = await client.post<ApiResponse<{ message: string }>>('/syslib/nfs/restart', {});
      return response.data;
    },

    listExports: async () => {
      const response = await client.get<ApiResponse<NFSExport[]>>('/syslib/nfs/exports');
      return response.data;
    },

    createExport: async (data: CreateNFSExportRequest) => {
      const response = await client.post<ApiResponse<{ message: string; path: string }>>('/syslib/nfs/exports', data);
      return response.data;
    },

    deleteExport: async (path: string) => {
      const response = await client.delete<ApiResponse<{ message: string }>>(`/syslib/nfs/exports?path=${encodeURIComponent(path)}`);
      return response.data;
    },
  },

  // Network Operations
  network: {
    createBond: async (data: CreateBondRequest) => {
      const response = await client.post<ApiResponse<{ message: string; bond: string }>>('/syslib/network/bond', data);
      return response.data;
    },

    deleteBond: async (name: string) => {
      const response = await client.delete<ApiResponse<{ message: string }>>(`/syslib/network/bond/${name}`);
      return response.data;
    },

    createVLAN: async (data: CreateVLANRequest) => {
      const response = await client.post<ApiResponse<{ message: string; vlan: string }>>('/syslib/network/vlan', data);
      return response.data;
    },

    deleteVLAN: async (parent: string, vlanId: number) => {
      const response = await client.delete<ApiResponse<{ message: string }>>(`/syslib/network/vlan/${parent}/${vlanId}`);
      return response.data;
    },
  },
};
