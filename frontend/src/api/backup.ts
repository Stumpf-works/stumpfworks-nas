import client, { ApiResponse } from './client';

// Types
export interface BackupJob {
  id: string;
  name: string;
  description: string;
  source: string;
  destination: string;
  type: string; // full, incremental, differential
  schedule: string; // cron expression
  enabled: boolean;
  retention: number; // days to keep backups
  compression: boolean;
  encryption: boolean;
  lastRun?: string;
  nextRun?: string;
  status: string; // idle, running, success, failed
  config?: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export interface BackupHistory {
  id: string;
  jobId: string;
  jobName: string;
  startTime: string;
  endTime?: string;
  status: string; // running, success, failed
  bytesBackup: number;
  filesBackup: number;
  duration: number; // seconds
  error?: string;
  backupPath: string;
}

export interface Snapshot {
  id: string;
  name: string;
  filesystem: string;
  createdAt: string;
  size: number;
  used: number;
  referenced: number;
  type: string; // zfs, btrfs, lvm
  description?: string;
}

export interface CreateSnapshotRequest {
  filesystem: string;
  name: string;
}

export interface RestoreSnapshotRequest {
  destination: string;
}

// API
export const backupApi = {
  // Backup Jobs
  async listJobs(): Promise<ApiResponse<BackupJob[]>> {
    const response = await client.get('/backups/jobs');
    return response.data;
  },

  async getJob(id: string): Promise<ApiResponse<BackupJob>> {
    const response = await client.get(`/backups/jobs/${id}`);
    return response.data;
  },

  async createJob(job: Partial<BackupJob>): Promise<ApiResponse<BackupJob>> {
    const response = await client.post('/backups/jobs', job);
    return response.data;
  },

  async updateJob(id: string, updates: Partial<BackupJob>): Promise<ApiResponse<any>> {
    const response = await client.put(`/backups/jobs/${id}`, updates);
    return response.data;
  },

  async deleteJob(id: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/backups/jobs/${id}`);
    return response.data;
  },

  async runJob(id: string): Promise<ApiResponse<BackupHistory>> {
    const response = await client.post(`/backups/jobs/${id}/run`);
    return response.data;
  },

  // Backup History
  async getHistory(jobId?: string, limit?: number): Promise<ApiResponse<BackupHistory[]>> {
    const params: any = {};
    if (jobId) params.jobId = jobId;
    if (limit) params.limit = limit;

    const response = await client.get('/backups/history', { params });
    return response.data;
  },

  // Snapshots
  async listSnapshots(): Promise<ApiResponse<Snapshot[]>> {
    const response = await client.get('/backups/snapshots');
    return response.data;
  },

  async createSnapshot(request: CreateSnapshotRequest): Promise<ApiResponse<Snapshot>> {
    const response = await client.post('/backups/snapshots', request);
    return response.data;
  },

  async deleteSnapshot(id: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/backups/snapshots/${id}`);
    return response.data;
  },

  async restoreSnapshot(
    id: string,
    request: RestoreSnapshotRequest
  ): Promise<ApiResponse<any>> {
    const response = await client.post(`/backups/snapshots/${id}/restore`, request);
    return response.data;
  },
};
