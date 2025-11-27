import client, { ApiResponse } from './client';

export interface FailedLoginAttempt {
  id: number;
  createdAt: string;
  username: string;
  ipAddress: string;
  userAgent?: string;
  reason: string;
  blocked: boolean;
  blockedAt?: string;
}

export interface IPBlock {
  id: number;
  createdAt: string;
  expiresAt: string;
  ipAddress: string;
  reason: string;
  attempts: number;
  isActive: boolean;
  isPermanent: boolean;
}

export interface FailedLoginStats {
  total_attempts: number;
  last_24h_attempts: number;
  blocked_ips: number;
  top_failed_usernames: Array<{
    username: string;
    count: number;
  }>;
}

export interface FailedLoginListResponse {
  attempts: FailedLoginAttempt[];
  total: number;
  limit: number;
  offset: number;
}

export const securityApi = {
  // Get failed login attempts with pagination
  getFailedLogins: async (limit?: number, offset?: number) => {
    const queryParams = new URLSearchParams();
    if (limit) queryParams.append('limit', limit.toString());
    if (offset) queryParams.append('offset', offset.toString());

    const response = await client.get<ApiResponse<FailedLoginListResponse>>(
      `/api/v1/security/failed-logins?${queryParams.toString()}`
    );
    return response.data;
  },

  // Get all blocked IPs
  getBlockedIPs: async () => {
    const response = await client.get<ApiResponse<IPBlock[]>>(
      '/security/blocked-ips'
    );
    return response.data;
  },

  // Unblock an IP address
  unblockIP: async (ipAddress: string) => {
    const response = await client.post<ApiResponse<{ message: string }>>(
      '/security/unblock-ip',
      { ipAddress }
    );
    return response.data;
  },

  // Get failed login statistics
  getStats: async () => {
    const response = await client.get<ApiResponse<FailedLoginStats>>(
      '/security/failed-logins/stats'
    );
    return response.data;
  },
};
