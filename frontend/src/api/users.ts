import client, { ApiResponse } from './client';
import { User } from './auth';

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  fullName?: string;
  role: 'admin' | 'user' | 'guest';
}

export interface UpdateUserRequest {
  email?: string;
  fullName?: string;
  role?: 'admin' | 'user' | 'guest';
  isActive?: boolean;
  password?: string;
}

export const usersApi = {
  list: async () => {
    const response = await client.get<ApiResponse<User[]>>('/users');
    return response.data;
  },

  get: async (id: number) => {
    const response = await client.get<ApiResponse<User>>(`/users/${id}`);
    return response.data;
  },

  create: async (data: CreateUserRequest) => {
    const response = await client.post<ApiResponse<User>>('/users', data);
    return response.data;
  },

  update: async (id: number, data: UpdateUserRequest) => {
    const response = await client.put<ApiResponse<User>>(`/users/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await client.delete<ApiResponse>(`/users/${id}`);
    return response.data;
  },
};
