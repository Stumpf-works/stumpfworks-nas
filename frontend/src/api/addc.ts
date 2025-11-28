// Revision: 2025-11-26 | Author: StumpfWorks AI | Version: 1.0.0
import client, { ApiResponse } from './client';

// ===== Types =====

// Domain Controller Configuration
export interface DCConfig {
  enabled: boolean;
  realm: string; // e.g., EXAMPLE.COM
  domain: string; // NetBIOS name, e.g., EXAMPLE
  server_role: string; // dc, member, standalone
  dns_backend: string; // SAMBA_INTERNAL, BIND9_DLZ, NONE
  dns_forwarder: string; // Forwarder IP
  function_level: string; // 2008_R2, 2012, 2012_R2, 2016
  host_ip: string; // Server IP
  sysvol_path: string; // Path to SYSVOL
  private_dir_path: string; // Path to private dir
}

// Domain Controller Status
export interface DCStatus {
  provisioned: boolean;
  config: DCConfig;
  service_status?: string;
  domain_info?: DomainInfo;
}

// Domain Information
export interface DomainInfo {
  Forest?: string;
  Domain?: string;
  Netbios?: string;
  DC?: string;
  [key: string]: any;
}

// Provision Options
export interface ProvisionOptions {
  realm: string; // e.g., EXAMPLE.COM
  domain: string; // NetBIOS domain name, e.g., EXAMPLE
  admin_password: string; // Administrator password
  dns_backend?: string; // SAMBA_INTERNAL, BIND9_DLZ, or NONE
  dns_forwarder?: string; // Optional DNS forwarder IP
  server_role?: string; // dc, member, standalone
  use_tls?: boolean; // Use LDAPS
  function_level?: string; // 2008_R2, 2012, 2012_R2, 2016
  host_ip?: string; // Server IP address
}

// Active Directory User (Extended)
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
  enabled: boolean;
  password_expired: boolean;
  member_of?: string[];
}

// Create User Request
export interface CreateUserRequest {
  user: ADDCUser;
  password: string;
}

// Active Directory Group
export interface ADGroup {
  name: string;
  description?: string;
  ou?: string;
  group_scope?: string; // Domain, Global, Universal
  group_type?: string; // Security, Distribution
  members?: string[];
}

// Active Directory Computer
export interface ADComputer {
  name: string;
  description?: string;
  ou?: string;
  ip?: string;
  enabled: boolean;
}

// Active Directory Organizational Unit
export interface ADOU {
  name: string;
  description?: string;
  parent_dn?: string; // e.g., "DC=example,DC=com"
}

// Active Directory Group Policy Object
export interface ADGPO {
  name: string;
  display_name?: string;
  description?: string;
}

// DNS Record
export interface ADDNSRecord {
  name: string;
  type: string; // A, AAAA, CNAME, MX, TXT, SRV
  value: string;
  ttl?: number;
}

// FSMO Roles
export interface FSMORoles {
  [roleName: string]: string;
}

// ===== API Client =====

export const addcApi = {
  // ===== Domain Controller Management =====

  // Get DC status
  getStatus: async () => {
    const response = await client.get<ApiResponse<DCStatus>>('/ad-dc/status');
    return response.data;
  },

  // Get DC configuration
  getConfig: async () => {
    const response = await client.get<ApiResponse<DCConfig>>('/ad-dc/config');
    return response.data;
  },

  // Update DC configuration
  updateConfig: async (config: DCConfig) => {
    const response = await client.put<ApiResponse<DCConfig>>('/ad-dc/config', config);
    return response.data;
  },

  // Provision new domain
  provisionDomain: async (options: ProvisionOptions) => {
    const response = await client.post<ApiResponse<{ message: string; realm: string; domain: string }>>(
      '/ad-dc/provision',
      options
    );
    return response.data;
  },

  // Demote domain controller
  demoteDomain: async () => {
    const response = await client.post<ApiResponse<{ message: string }>>('/ad-dc/demote', {});
    return response.data;
  },

  // Get domain information
  getDomainInfo: async () => {
    const response = await client.get<ApiResponse<DomainInfo>>('/ad-dc/info');
    return response.data;
  },

  // Get domain functional level
  getDomainLevel: async () => {
    const response = await client.get<ApiResponse<{ level: string }>>('/ad-dc/level');
    return response.data;
  },

  // Raise domain functional level
  raiseDomainLevel: async (level: string) => {
    const response = await client.post<ApiResponse<{ message: string; level: string }>>(
      '/ad-dc/level/raise',
      { level }
    );
    return response.data;
  },

  // Restart Samba AD DC service
  restartService: async () => {
    const response = await client.post<ApiResponse<{ message: string }>>('/ad-dc/service/restart', {});
    return response.data;
  },

  // ===== User Management =====

  // List all users
  listUsers: async () => {
    const response = await client.get<ApiResponse<string[]>>('/ad-dc/users');
    return response.data;
  },

  // Create new user
  createUser: async (request: CreateUserRequest) => {
    const response = await client.post<ApiResponse<{ message: string; username: string }>>(
      '/ad-dc/users',
      request
    );
    return response.data;
  },

  // Delete user
  deleteUser: async (username: string) => {
    const response = await client.delete<ApiResponse<{ message: string; username: string }>>(
      `/ad-dc/users/${username}`
    );
    return response.data;
  },

  // Enable user
  enableUser: async (username: string) => {
    const response = await client.post<ApiResponse<{ message: string; username: string }>>(
      `/ad-dc/users/${username}/enable`,
      {}
    );
    return response.data;
  },

  // Disable user
  disableUser: async (username: string) => {
    const response = await client.post<ApiResponse<{ message: string; username: string }>>(
      `/ad-dc/users/${username}/disable`,
      {}
    );
    return response.data;
  },

  // Set user password
  setUserPassword: async (username: string, password: string) => {
    const response = await client.post<ApiResponse<{ message: string; username: string }>>(
      `/ad-dc/users/${username}/password`,
      { password }
    );
    return response.data;
  },

  // Set user expiry
  setUserExpiry: async (username: string, days: number) => {
    const response = await client.post<ApiResponse<{ message: string; username: string; days: number }>>(
      `/ad-dc/users/${username}/expiry`,
      { days }
    );
    return response.data;
  },

  // ===== Group Management =====

  // List all groups
  listGroups: async () => {
    const response = await client.get<ApiResponse<string[]>>('/ad-dc/groups');
    return response.data;
  },

  // Create new group
  createGroup: async (group: ADGroup) => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>(
      '/ad-dc/groups',
      group
    );
    return response.data;
  },

  // Delete group
  deleteGroup: async (name: string) => {
    const response = await client.delete<ApiResponse<{ message: string; name: string }>>(
      `/ad-dc/groups/${name}`
    );
    return response.data;
  },

  // List group members
  listGroupMembers: async (name: string) => {
    const response = await client.get<ApiResponse<string[]>>(`/ad-dc/groups/${name}/members`);
    return response.data;
  },

  // Add member to group
  addGroupMember: async (name: string, username: string) => {
    const response = await client.post<ApiResponse<{ message: string; group: string; username: string }>>(
      `/ad-dc/groups/${name}/members`,
      { username }
    );
    return response.data;
  },

  // Remove member from group
  removeGroupMember: async (name: string, username: string) => {
    const response = await client.delete<ApiResponse<{ message: string; group: string; username: string }>>(
      `/ad-dc/groups/${name}/members/${username}`
    );
    return response.data;
  },

  // ===== Computer Management =====

  // List all computers
  listComputers: async () => {
    const response = await client.get<ApiResponse<string[]>>('/ad-dc/computers');
    return response.data;
  },

  // Create new computer
  createComputer: async (computer: ADComputer) => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>(
      '/ad-dc/computers',
      computer
    );
    return response.data;
  },

  // Delete computer
  deleteComputer: async (name: string) => {
    const response = await client.delete<ApiResponse<{ message: string; name: string }>>(
      `/ad-dc/computers/${name}`
    );
    return response.data;
  },

  // ===== Organizational Unit Management =====

  // List all OUs
  listOUs: async () => {
    const response = await client.get<ApiResponse<string[]>>('/ad-dc/ou');
    return response.data;
  },

  // Create new OU
  createOU: async (ou: ADOU) => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>('/ad-dc/ou', ou);
    return response.data;
  },

  // Delete OU
  deleteOU: async (dn: string) => {
    const response = await client.delete<ApiResponse<{ message: string; dn: string }>>('/ad-dc/ou', {
      data: { dn },
    });
    return response.data;
  },

  // ===== Group Policy Management =====

  // List all GPOs
  listGPOs: async () => {
    const response = await client.get<ApiResponse<string[]>>('/ad-dc/gpo');
    return response.data;
  },

  // Create new GPO
  createGPO: async (gpo: ADGPO) => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>('/ad-dc/gpo', gpo);
    return response.data;
  },

  // Delete GPO
  deleteGPO: async (name: string) => {
    const response = await client.delete<ApiResponse<{ message: string; name: string }>>(`/ad-dc/gpo/${name}`);
    return response.data;
  },

  // Link GPO to OU
  linkGPO: async (name: string, ouDN: string) => {
    const response = await client.post<ApiResponse<{ message: string; gpo: string; ou: string }>>(
      `/ad-dc/gpo/${name}/link`,
      { ou_dn: ouDN }
    );
    return response.data;
  },

  // Unlink GPO from OU
  unlinkGPO: async (name: string, ouDN: string) => {
    const response = await client.post<ApiResponse<{ message: string; gpo: string; ou: string }>>(
      `/ad-dc/gpo/${name}/unlink`,
      { ou_dn: ouDN }
    );
    return response.data;
  },

  // ===== DNS Management =====

  // List DNS zones
  listDNSZones: async () => {
    const response = await client.get<ApiResponse<string[]>>('/ad-dc/dns/zones');
    return response.data;
  },

  // Create DNS zone
  createDNSZone: async (zoneName: string) => {
    const response = await client.post<ApiResponse<{ message: string; zone: string }>>(
      '/ad-dc/dns/zones',
      { zone_name: zoneName }
    );
    return response.data;
  },

  // Delete DNS zone
  deleteDNSZone: async (zone: string) => {
    const response = await client.delete<ApiResponse<{ message: string; zone: string }>>(
      `/ad-dc/dns/zones/${zone}`
    );
    return response.data;
  },

  // List DNS records in zone
  listDNSRecords: async (zone: string) => {
    const response = await client.get<ApiResponse<string[]>>(`/ad-dc/dns/zones/${zone}/records`);
    return response.data;
  },

  // Add DNS record
  addDNSRecord: async (zone: string, record: ADDNSRecord) => {
    const response = await client.post<ApiResponse<{ message: string; zone: string; record: string }>>(
      `/ad-dc/dns/zones/${zone}/records`,
      record
    );
    return response.data;
  },

  // Delete DNS record
  deleteDNSRecord: async (zone: string, record: string, recordType: string, value: string) => {
    const response = await client.delete<ApiResponse<{ message: string; zone: string; record: string }>>(
      `/ad-dc/dns/zones/${zone}/records/${record}`,
      {
        data: { type: recordType, value },
      }
    );
    return response.data;
  },

  // ===== FSMO Roles Management =====

  // Show FSMO roles
  showFSMORoles: async () => {
    const response = await client.get<ApiResponse<FSMORoles>>('/ad-dc/fsmo');
    return response.data;
  },

  // Transfer FSMO roles
  transferFSMORoles: async (role: string, targetDC: string) => {
    const response = await client.post<ApiResponse<{ message: string; role: string; target_dc: string }>>(
      '/ad-dc/fsmo/transfer',
      { role, target_dc: targetDC }
    );
    return response.data;
  },

  // Seize FSMO roles
  seizeFSMORoles: async (role: string) => {
    const response = await client.post<ApiResponse<{ message: string; role: string }>>(
      '/ad-dc/fsmo/seize',
      { role }
    );
    return response.data;
  },

  // ===== Utility Functions =====

  // Test configuration
  testConfiguration: async () => {
    const response = await client.post<ApiResponse<{ message: string }>>('/ad-dc/test-config', {});
    return response.data;
  },

  // Show database check
  showDBCheck: async () => {
    const response = await client.get<ApiResponse<{ result: string }>>('/ad-dc/dbcheck');
    return response.data;
  },

  // Perform online backup
  backupOnline: async (targetDir: string) => {
    const response = await client.post<ApiResponse<{ message: string; target_dir: string }>>(
      '/ad-dc/backup',
      { target_dir: targetDir }
    );
    return response.data;
  },
};
