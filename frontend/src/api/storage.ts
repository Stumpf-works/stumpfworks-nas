import client, { ApiResponse } from './client';

// ===== Types =====

export interface Disk {
  name: string;
  path: string;
  label?: string; // User-defined friendly name
  model: string;
  serial: string;
  size: number;
  type: 'hdd' | 'ssd' | 'nvme' | 'usb';
  status: 'healthy' | 'warning' | 'critical' | 'failed' | 'unknown';
  temperature: number;
  isSystem: boolean;
  isRemovable: boolean;
  partitions: Partition[];
  smartEnabled: boolean;
  smart?: SMARTData;
}

export interface Partition {
  name: string;
  path: string;
  size: number;
  used: number;
  filesystem: string;
  mountPoint: string;
  label: string;
  uuid: string;
  isMounted: boolean;
}

export interface SMARTData {
  healthy: boolean;
  temperature: number;
  powerOnHours: number;
  powerCycleCount: number;
  reallocatedSectors: number;
  pendingSectors: number;
  uncorrectableErrors: number;
  crcErrors: number;
  percentLifeUsed: number;
  lastUpdated: string;
}

export interface Volume {
  id: string;
  name: string;
  type: 'single' | 'raid0' | 'raid1' | 'raid5' | 'raid6' | 'raid10' | 'lvm' | 'zfs' | 'btrfs';
  status: 'online' | 'degraded' | 'offline' | 'rebuilding' | 'failed';
  size: number;
  used: number;
  available: number;
  filesystem: string;
  mountPoint: string;
  disks: string[];
  raidLevel?: string;
  health: number;
  createdAt: string;
  snapshots?: Snapshot[];
}

export interface Snapshot {
  id: string;
  volumeId: string;
  name: string;
  size: number;
  createdAt: string;
}

export interface Share {
  id: string;
  name: string;
  path: string;
  volumeId?: string; // Optional - linked volume
  type: 'smb' | 'nfs' | 'ftp';
  description: string;
  enabled: boolean;
  readOnly: boolean;
  browseable: boolean;
  guestOk: boolean;
  validUsers?: string[];
  validGroups?: string[];
  createdAt: string;
  updatedAt: string;
}

export interface StorageStats {
  totalDisks: number;
  totalCapacity: number;
  usedCapacity: number;
  availableCapacity: number;
  totalVolumes: number;
  totalShares: number;
  healthyDisks: number;
  warningDisks: number;
  criticalDisks: number;
}

export interface DiskIOStats {
  diskName: string;
  readBytes: number;
  writeBytes: number;
  readOps: number;
  writeOps: number;
  readLatency: number;
  writeLatency: number;
  utilization: number;
  timestamp: string;
}

export interface DiskHealth {
  diskName: string;
  status: 'healthy' | 'warning' | 'critical' | 'failed' | 'unknown';
  issues: string[];
  temperature: number;
  score: number;
}

export interface CreateVolumeRequest {
  name: string;
  type: 'single' | 'raid0' | 'raid1' | 'raid5' | 'raid6' | 'raid10' | 'lvm';
  disks: string[];
  filesystem: 'ext4' | 'xfs' | 'btrfs' | 'zfs';
  mountPoint: string;
  raidLevel?: string;
}

export interface CreateShareRequest {
  name: string;
  volumeId?: string; // Optional - select from managed volumes
  path?: string;     // Optional - manual path (used if volumeId not provided)
  type: 'smb' | 'nfs' | 'ftp';
  description: string;
  readOnly: boolean;
  browseable: boolean;
  guestOk: boolean;
  validUsers?: string[];
  validGroups?: string[];
}

export interface FormatDiskRequest {
  disk: string;
  filesystem: 'ext4' | 'xfs' | 'btrfs';
  label?: string;
  force: boolean;
}

// ===== API Client =====

export const storageApi = {
  // Statistics
  getStats: async () => {
    const response = await client.get<ApiResponse<StorageStats>>('/storage/stats');
    return response.data;
  },

  getHealth: async () => {
    const response = await client.get<ApiResponse<DiskHealth[]>>('/storage/health');
    return response.data;
  },

  getIOStats: async () => {
    const response = await client.get<ApiResponse<DiskIOStats[]>>('/storage/io');
    return response.data;
  },

  getIOMonitor: async () => {
    const response = await client.get<ApiResponse<DiskIOStats[]>>('/storage/io/monitor');
    return response.data;
  },

  // Disks
  listDisks: async () => {
    const response = await client.get<ApiResponse<Disk[]>>('/storage/disks');
    return response.data;
  },

  getDisk: async (name: string) => {
    const response = await client.get<ApiResponse<Disk>>(`/storage/disks/${name}`);
    return response.data;
  },

  getDiskSMART: async (name: string) => {
    const response = await client.get<ApiResponse<SMARTData>>(`/storage/disks/${name}/smart`);
    return response.data;
  },

  getDiskHealth: async (name: string) => {
    const response = await client.get<ApiResponse<DiskHealth>>(`/storage/disks/${name}/health`);
    return response.data;
  },

  getDiskIO: async (name: string) => {
    const response = await client.get<ApiResponse<DiskIOStats>>(`/storage/disks/${name}/io`);
    return response.data;
  },

  formatDisk: async (data: FormatDiskRequest) => {
    const response = await client.post<ApiResponse<{ message: string }>>('/storage/disks/format', data);
    return response.data;
  },

  setDiskLabel: async (diskName: string, label: string) => {
    const response = await client.put<ApiResponse<{ message: string; label: string }>>(`/storage/disks/${diskName}/label`, { label });
    return response.data;
  },

  // Volumes
  listVolumes: async () => {
    const response = await client.get<ApiResponse<Volume[]>>('/storage/volumes');
    return response.data;
  },

  getVolume: async (id: string) => {
    const response = await client.get<ApiResponse<Volume>>(`/storage/volumes/${id}`);
    return response.data;
  },

  createVolume: async (data: CreateVolumeRequest) => {
    const response = await client.post<ApiResponse<Volume>>('/storage/volumes', data);
    return response.data;
  },

  deleteVolume: async (id: string) => {
    const response = await client.delete<ApiResponse<{ message: string }>>(`/storage/volumes/${id}`);
    return response.data;
  },

  // Shares
  listShares: async () => {
    const response = await client.get<ApiResponse<Share[]>>('/storage/shares');
    return response.data;
  },

  getShare: async (id: string) => {
    const response = await client.get<ApiResponse<Share>>(`/storage/shares/${id}`);
    return response.data;
  },

  createShare: async (data: CreateShareRequest) => {
    const response = await client.post<ApiResponse<Share>>('/storage/shares', data);
    return response.data;
  },

  updateShare: async (id: string, data: CreateShareRequest) => {
    const response = await client.put<ApiResponse<Share>>(`/storage/shares/${id}`, data);
    return response.data;
  },

  deleteShare: async (id: string) => {
    const response = await client.delete<ApiResponse<{ message: string }>>(`/storage/shares/${id}`);
    return response.data;
  },

  enableShare: async (id: string) => {
    const response = await client.post<ApiResponse<{ message: string }>>(`/storage/shares/${id}/enable`);
    return response.data;
  },

  disableShare: async (id: string) => {
    const response = await client.post<ApiResponse<{ message: string }>>(`/storage/shares/${id}/disable`);
    return response.data;
  },
};
