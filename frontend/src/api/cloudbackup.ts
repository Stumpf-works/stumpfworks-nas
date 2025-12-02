import client, { ApiResponse } from './client';

// Types
export interface CloudProvider {
  id: number;
  created_at: string;
  updated_at: string;
  name: string;
  type: string; // s3, b2, gdrive, dropbox, onedrive, azureblob, sftp
  description: string;
  config: string; // JSON-encoded config
  enabled: boolean;
  test_status: string; // untested, success, failed
  tested_at?: string;
}

export interface CloudSyncJob {
  id: number;
  created_at: string;
  updated_at: string;
  name: string;
  description: string;
  enabled: boolean;
  provider_id: number;
  provider?: CloudProvider;
  direction: string; // upload, download, sync
  local_path: string;
  remote_path: string;
  schedule: string; // cron expression
  schedule_enabled: boolean;
  bandwidth_limit: string; // e.g., "10M", "1G"
  encryption_enabled: boolean;
  delete_after_upload: boolean;
  retention: number; // days to keep backups (0 = forever)
  filters: string; // JSON array of include/exclude patterns
  last_run_at?: string;
  next_run_at?: string;
  last_status: string; // idle, running, success, failed
  failure_count: number;
}

export interface CloudSyncLog {
  id: number;
  created_at: string;
  updated_at: string;
  job_id: number;
  job?: CloudSyncJob;
  job_name: string;
  started_at: string;
  completed_at?: string;
  status: string; // running, success, failed, cancelled
  direction: string; // upload, download, sync
  bytes_transferred: number;
  files_transferred: number;
  files_deleted: number;
  files_failed: number;
  duration: number; // seconds
  error_message?: string;
  output?: string;
  triggered_by: string; // manual, schedule, api
}

export interface ProviderTypeConfig {
  [key: string]: {
    [field: string]: string;
  };
}

// API
export const cloudBackupApi = {
  // Providers
  async listProviders(): Promise<ApiResponse<CloudProvider[]>> {
    const response = await client.get('/cloudbackup/providers');
    return response.data;
  },

  async getProvider(id: number): Promise<ApiResponse<CloudProvider>> {
    const response = await client.get(`/cloudbackup/providers/${id}`);
    return response.data;
  },

  async createProvider(provider: Partial<CloudProvider>): Promise<ApiResponse<CloudProvider>> {
    const response = await client.post('/cloudbackup/providers', provider);
    return response.data;
  },

  async updateProvider(id: number, updates: Partial<CloudProvider>): Promise<ApiResponse<any>> {
    const response = await client.put(`/cloudbackup/providers/${id}`, updates);
    return response.data;
  },

  async deleteProvider(id: number): Promise<ApiResponse<any>> {
    const response = await client.delete(`/cloudbackup/providers/${id}`);
    return response.data;
  },

  async testProvider(id: number): Promise<ApiResponse<any>> {
    const response = await client.post(`/cloudbackup/providers/${id}/test`);
    return response.data;
  },

  async getProviderTypes(): Promise<ApiResponse<ProviderTypeConfig>> {
    const response = await client.get('/cloudbackup/providers/types');
    return response.data;
  },

  // Jobs
  async listJobs(): Promise<ApiResponse<CloudSyncJob[]>> {
    const response = await client.get('/cloudbackup/jobs');
    return response.data;
  },

  async getJob(id: number): Promise<ApiResponse<CloudSyncJob>> {
    const response = await client.get(`/cloudbackup/jobs/${id}`);
    return response.data;
  },

  async createJob(job: Partial<CloudSyncJob>): Promise<ApiResponse<CloudSyncJob>> {
    const response = await client.post('/cloudbackup/jobs', job);
    return response.data;
  },

  async updateJob(id: number, updates: Partial<CloudSyncJob>): Promise<ApiResponse<any>> {
    const response = await client.put(`/cloudbackup/jobs/${id}`, updates);
    return response.data;
  },

  async deleteJob(id: number): Promise<ApiResponse<any>> {
    const response = await client.delete(`/cloudbackup/jobs/${id}`);
    return response.data;
  },

  async runJob(id: number): Promise<ApiResponse<CloudSyncLog>> {
    const response = await client.post(`/cloudbackup/jobs/${id}/run`);
    return response.data;
  },

  // Logs
  async getLogs(jobId?: number, limit?: number, offset?: number): Promise<ApiResponse<{
    logs: CloudSyncLog[];
    total: number;
    limit: number;
    offset: number;
  }>> {
    const params: any = {};
    if (jobId) params.jobId = jobId;
    if (limit) params.limit = limit;
    if (offset) params.offset = offset;

    const response = await client.get('/cloudbackup/logs', { params });
    return response.data;
  },
};
