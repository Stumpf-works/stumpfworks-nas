import client, { ApiResponse } from './client';
import { User } from './auth';

export interface UserGroup {
  id: number;
  name: string;
  description: string;
  isSystem: boolean;
  memberCount: number;
  members?: MemberInfo[];
  createdAt: string;
  updatedAt: string;
}

export interface MemberInfo {
  id: number;
  username: string;
  fullName: string;
  email: string;
}

export interface CreateGroupRequest {
  name: string;
  description?: string;
}

export interface UpdateGroupRequest {
  name?: string;
  description?: string;
}

export interface AddMemberRequest {
  userId: number;
}

export const groupsApi = {
  list: async () => {
    const response = await client.get<ApiResponse<UserGroup[]>>('/groups');
    return response.data;
  },

  get: async (id: number) => {
    const response = await client.get<ApiResponse<UserGroup>>(`/groups/${id}`);
    return response.data;
  },

  create: async (data: CreateGroupRequest) => {
    const response = await client.post<ApiResponse<UserGroup>>('/groups', data);
    return response.data;
  },

  update: async (id: number, data: UpdateGroupRequest) => {
    const response = await client.put<ApiResponse<UserGroup>>(`/groups/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await client.delete<ApiResponse>(`/groups/${id}`);
    return response.data;
  },

  addMember: async (groupId: number, userId: number) => {
    const response = await client.post<ApiResponse>(
      `/groups/${groupId}/members`,
      { userId }
    );
    return response.data;
  },

  removeMember: async (groupId: number, userId: number) => {
    const response = await client.delete<ApiResponse>(
      `/groups/${groupId}/members/${userId}`
    );
    return response.data;
  },

  getMembers: async (groupId: number) => {
    const response = await client.get<ApiResponse<User[]>>(
      `/groups/${groupId}/members`
    );
    return response.data;
  },
};
