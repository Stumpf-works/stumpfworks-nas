import client, { ApiResponse } from './client';

export interface ACLEntry {
  type: string;        // user, group, mask, other
  name: string;        // username or groupname (empty for owner/other)
  permissions: string; // rwx format (e.g., "rwx", "r-x", "---")
}

export interface ACLInfo {
  path: string;
  entries: ACLEntry[];
}

export interface SetACLRequest {
  path: string;
  entries: ACLEntry[];
}

export interface RemoveACLRequest {
  path: string;
  type: string;
  name: string;
}

export interface SetDefaultACLRequest {
  dir_path: string;
  entries: ACLEntry[];
}

export interface ApplyRecursiveRequest {
  dir_path: string;
  entries: ACLEntry[];
}

export const filesystemACLApi = {
  // Get ACL entries for a file or directory
  getACL: async (path: string) => {
    const response = await client.get<ApiResponse<ACLInfo>>(
      `/filesystem/acl?path=${encodeURIComponent(path)}`
    );
    return response.data;
  },

  // Set ACL entries on a file or directory
  setACL: async (request: SetACLRequest) => {
    const response = await client.post<ApiResponse<{ message: string; path: string }>>(
      '/filesystem/acl',
      request
    );
    return response.data;
  },

  // Remove a specific ACL entry
  removeACL: async (request: RemoveACLRequest) => {
    const response = await client.delete<ApiResponse<{ message: string; path: string }>>(
      '/filesystem/acl',
      { data: request }
    );
    return response.data;
  },

  // Set default ACL entries for a directory
  setDefaultACL: async (request: SetDefaultACLRequest) => {
    const response = await client.post<ApiResponse<{ message: string; path: string }>>(
      '/filesystem/acl/default',
      request
    );
    return response.data;
  },

  // Apply ACL entries recursively to a directory
  applyRecursive: async (request: ApplyRecursiveRequest) => {
    const response = await client.post<ApiResponse<{ message: string; path: string }>>(
      '/filesystem/acl/recursive',
      request
    );
    return response.data;
  },

  // Remove all ACL entries from a file or directory
  removeAllACLs: async (path: string) => {
    const response = await client.delete<ApiResponse<{ message: string; path: string }>>(
      '/filesystem/acl/all',
      { data: { path } }
    );
    return response.data;
  },
};
