import client, { ApiResponse } from './client';

export interface AuditLog {
  id: number;
  createdAt: string;
  userId?: number;
  username: string;
  action: string;
  resource?: string;
  status: 'success' | 'failure' | 'error';
  severity: 'info' | 'warning' | 'critical';
  ipAddress?: string;
  userAgent?: string;
  details?: string;
  message?: string;
}

export interface AuditLogQueryParams {
  userId?: number;
  username?: string;
  action?: string;
  status?: string;
  severity?: string;
  startDate?: string;
  endDate?: string;
  limit?: number;
  offset?: number;
}

export interface AuditLogListResponse {
  logs: AuditLog[];
  total: number;
  limit: number;
  offset: number;
}

export interface AuditStats {
  total: number;
  last_24h: number;
  by_severity: Array<{
    severity: string;
    count: number;
  }>;
  top_actions: Array<{
    action: string;
    count: number;
  }>;
}

export const auditApi = {
  // List audit logs with filters and pagination
  listLogs: async (params?: AuditLogQueryParams) => {
    const queryParams = new URLSearchParams();

    if (params?.userId) queryParams.append('userId', params.userId.toString());
    if (params?.username) queryParams.append('username', params.username);
    if (params?.action) queryParams.append('action', params.action);
    if (params?.status) queryParams.append('status', params.status);
    if (params?.severity) queryParams.append('severity', params.severity);
    if (params?.startDate) queryParams.append('startDate', params.startDate);
    if (params?.endDate) queryParams.append('endDate', params.endDate);
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const response = await client.get<ApiResponse<AuditLogListResponse>>(
      `/api/v1/audit/logs?${queryParams.toString()}`
    );
    return response.data;
  },

  // Get a specific audit log by ID
  getLog: async (id: number) => {
    const response = await client.get<ApiResponse<AuditLog>>(
      `/api/v1/audit/logs/${id}`
    );
    return response.data;
  },

  // Get recent audit logs
  getRecent: async (limit?: number) => {
    const queryParams = limit ? `?limit=${limit}` : '';
    const response = await client.get<ApiResponse<AuditLog[]>>(
      `/api/v1/audit/logs/recent${queryParams}`
    );
    return response.data;
  },

  // Get audit log statistics
  getStats: async () => {
    const response = await client.get<ApiResponse<AuditStats>>(
      '/api/v1/audit/stats'
    );
    return response.data;
  },
};
