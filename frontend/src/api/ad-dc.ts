import { ApiResponse } from './client';

// ===== Domain Controller Types =====

export interface DCConfig {
  enabled: boolean;
  realm: string;
  domain: string;
  server_role: string;
  dns_backend: string;
  dns_forwarder: string;
  function_level: string;
  host_ip: string;
  sysvol_path: string;
  private_dir_path: string;
}

export interface DCStatus {
  provisioned: boolean;
  config: DCConfig;
  service_status?: string;
  domain_info?: Record<string, any>;
}

export interface ProvisionOptions {
  realm: string;
  domain: string;
  admin_password: string;
  dns_backend?: string;
  dns_forwarder?: string;
  server_role?: string;
  use_tls?: boolean;
  function_level?: string;
  host_ip?: string;
}

// ===== User Management Types =====

export interface ADDCUser {
  username: string;
  given_name?: string;
  surname?: string;
  display_name?: string;
  email?: string;
  description?: string;
  department?: string;
  company?: string;
  title?: string;
  telephone?: string;
  ou?: string;
  enabled?: boolean;
  password_expired?: boolean;
  member_of?: string[];
}

export interface CreateUserRequest {
  user: ADDCUser;
  password: string;
}

// ===== Group Management Types =====

export interface ADGroup {
  name: string;
  description?: string;
  ou?: string;
  group_scope?: string;
  group_type?: string;
  members?: string[];
}

// ===== Computer Management Types =====

export interface ADComputer {
  name: string;
  description?: string;
  ou?: string;
  ip?: string;
  enabled?: boolean;
}

// ===== OU Management Types =====

export interface ADOU {
  name: string;
  description?: string;
  parent_dn?: string;
}

// ===== GPO Management Types =====

export interface ADGPO {
  name: string;
  display_name?: string;
  description?: string;
}

// ===== DNS Management Types =====

export interface ADDNSRecord {
  name: string;
  type: string; // A, AAAA, CNAME, MX, TXT, SRV
  value: string;
  ttl?: number;
}

// ===== API Client =====

export const adDCApi = {
  // ===== Domain Controller Management =====

  getStatus: async (): Promise<ApiResponse<DCStatus>> => {
    const response = await fetch('/api/v1/ad-dc/status');
    return response.json();
  },

  getConfig: async (): Promise<ApiResponse<DCConfig>> => {
    const response = await fetch('/api/v1/ad-dc/config');
    return response.json();
  },

  updateConfig: async (config: DCConfig): Promise<ApiResponse<DCConfig>> => {
    const response = await fetch('/api/v1/ad-dc/config', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config),
    });
    return response.json();
  },

  provisionDomain: async (options: ProvisionOptions): Promise<ApiResponse<{ message: string; realm: string; domain: string }>> => {
    const response = await fetch('/api/v1/ad-dc/provision', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(options),
    });
    return response.json();
  },

  demoteDomain: async (): Promise<ApiResponse<{ message: string }>> => {
    const response = await fetch('/api/v1/ad-dc/demote', {
      method: 'POST',
    });
    return response.json();
  },

  getDomainInfo: async (): Promise<ApiResponse<Record<string, any>>> => {
    const response = await fetch('/api/v1/ad-dc/info');
    return response.json();
  },

  getDomainLevel: async (): Promise<ApiResponse<{ level: string }>> => {
    const response = await fetch('/api/v1/ad-dc/level');
    return response.json();
  },

  raiseDomainLevel: async (level: string): Promise<ApiResponse<{ message: string; level: string }>> => {
    const response = await fetch('/api/v1/ad-dc/level/raise', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ level }),
    });
    return response.json();
  },

  restartService: async (): Promise<ApiResponse<{ message: string }>> => {
    const response = await fetch('/api/v1/ad-dc/service/restart', {
      method: 'POST',
    });
    return response.json();
  },

  // ===== User Management =====

  listUsers: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ad-dc/users');
    return response.json();
  },

  createUser: async (request: CreateUserRequest): Promise<ApiResponse<{ message: string; username: string }>> => {
    const response = await fetch('/api/v1/ad-dc/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
    return response.json();
  },

  deleteUser: async (username: string): Promise<ApiResponse<{ message: string; username: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/users/${encodeURIComponent(username)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  enableUser: async (username: string): Promise<ApiResponse<{ message: string; username: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/users/${encodeURIComponent(username)}/enable`, {
      method: 'POST',
    });
    return response.json();
  },

  disableUser: async (username: string): Promise<ApiResponse<{ message: string; username: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/users/${encodeURIComponent(username)}/disable`, {
      method: 'POST',
    });
    return response.json();
  },

  setUserPassword: async (username: string, password: string): Promise<ApiResponse<{ message: string; username: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/users/${encodeURIComponent(username)}/password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password }),
    });
    return response.json();
  },

  setUserExpiry: async (username: string, days: number): Promise<ApiResponse<{ message: string; username: string; days: number }>> => {
    const response = await fetch(`/api/v1/ad-dc/users/${encodeURIComponent(username)}/expiry`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ days }),
    });
    return response.json();
  },

  // ===== Group Management =====

  listGroups: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ad-dc/groups');
    return response.json();
  },

  createGroup: async (group: ADGroup): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch('/api/v1/ad-dc/groups', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(group),
    });
    return response.json();
  },

  deleteGroup: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/groups/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  listGroupMembers: async (name: string): Promise<ApiResponse<string[]>> => {
    const response = await fetch(`/api/v1/ad-dc/groups/${encodeURIComponent(name)}/members`);
    return response.json();
  },

  addGroupMember: async (groupName: string, username: string): Promise<ApiResponse<{ message: string; group: string; username: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/groups/${encodeURIComponent(groupName)}/members`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username }),
    });
    return response.json();
  },

  removeGroupMember: async (groupName: string, username: string): Promise<ApiResponse<{ message: string; group: string; username: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/groups/${encodeURIComponent(groupName)}/members/${encodeURIComponent(username)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  // ===== Computer Management =====

  listComputers: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ad-dc/computers');
    return response.json();
  },

  createComputer: async (computer: ADComputer): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch('/api/v1/ad-dc/computers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(computer),
    });
    return response.json();
  },

  deleteComputer: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/computers/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  // ===== OU Management =====

  listOUs: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ad-dc/ou');
    return response.json();
  },

  createOU: async (ou: ADOU): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch('/api/v1/ad-dc/ou', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(ou),
    });
    return response.json();
  },

  deleteOU: async (dn: string): Promise<ApiResponse<{ message: string; dn: string }>> => {
    const response = await fetch('/api/v1/ad-dc/ou', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ dn }),
    });
    return response.json();
  },

  // ===== GPO Management =====

  listGPOs: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ad-dc/gpo');
    return response.json();
  },

  createGPO: async (gpo: ADGPO): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch('/api/v1/ad-dc/gpo', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(gpo),
    });
    return response.json();
  },

  deleteGPO: async (name: string): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/gpo/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  linkGPO: async (gpoName: string, ouDN: string): Promise<ApiResponse<{ message: string; gpo: string; ou: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/gpo/${encodeURIComponent(gpoName)}/link`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ou_dn: ouDN }),
    });
    return response.json();
  },

  unlinkGPO: async (gpoName: string, ouDN: string): Promise<ApiResponse<{ message: string; gpo: string; ou: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/gpo/${encodeURIComponent(gpoName)}/unlink`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ou_dn: ouDN }),
    });
    return response.json();
  },

  // ===== DNS Management =====

  listDNSZones: async (): Promise<ApiResponse<string[]>> => {
    const response = await fetch('/api/v1/ad-dc/dns/zones');
    return response.json();
  },

  createDNSZone: async (zoneName: string): Promise<ApiResponse<{ message: string; zone: string }>> => {
    const response = await fetch('/api/v1/ad-dc/dns/zones', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ zone_name: zoneName }),
    });
    return response.json();
  },

  deleteDNSZone: async (zoneName: string): Promise<ApiResponse<{ message: string; zone: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/dns/zones/${encodeURIComponent(zoneName)}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  listDNSRecords: async (zoneName: string): Promise<ApiResponse<string[]>> => {
    const response = await fetch(`/api/v1/ad-dc/dns/zones/${encodeURIComponent(zoneName)}/records`);
    return response.json();
  },

  addDNSRecord: async (zoneName: string, record: ADDNSRecord): Promise<ApiResponse<{ message: string; zone: string; record: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/dns/zones/${encodeURIComponent(zoneName)}/records`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(record),
    });
    return response.json();
  },

  deleteDNSRecord: async (zoneName: string, recordName: string, recordType: string, value: string): Promise<ApiResponse<{ message: string; zone: string; record: string }>> => {
    const response = await fetch(`/api/v1/ad-dc/dns/zones/${encodeURIComponent(zoneName)}/records/${encodeURIComponent(recordName)}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: recordType, value }),
    });
    return response.json();
  },

  // ===== FSMO Roles =====

  showFSMORoles: async (): Promise<ApiResponse<Record<string, string>>> => {
    const response = await fetch('/api/v1/ad-dc/fsmo');
    return response.json();
  },

  transferFSMORoles: async (role: string, targetDC: string): Promise<ApiResponse<{ message: string; role: string; target_dc: string }>> => {
    const response = await fetch('/api/v1/ad-dc/fsmo/transfer', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ role, target_dc: targetDC }),
    });
    return response.json();
  },

  seizeFSMORoles: async (role: string): Promise<ApiResponse<{ message: string; role: string }>> => {
    const response = await fetch('/api/v1/ad-dc/fsmo/seize', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ role }),
    });
    return response.json();
  },

  // ===== Utility Functions =====

  testConfiguration: async (): Promise<ApiResponse<{ message: string }>> => {
    const response = await fetch('/api/v1/ad-dc/test-config', {
      method: 'POST',
    });
    return response.json();
  },

  showDBCheck: async (): Promise<ApiResponse<{ result: string }>> => {
    const response = await fetch('/api/v1/ad-dc/dbcheck');
    return response.json();
  },

  backupOnline: async (targetDir: string): Promise<ApiResponse<{ message: string; target_dir: string }>> => {
    const response = await fetch('/api/v1/ad-dc/backup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_dir: targetDir }),
    });
    return response.json();
  },
};
