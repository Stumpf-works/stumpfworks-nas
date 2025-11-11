import client, { ApiResponse } from './client';

export interface User {
  id: number;
  username: string;
  email: string;
  fullName: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

export const authApi = {
  login: async (credentials: LoginRequest) => {
    const response = await client.post<ApiResponse<LoginResponse>>('/auth/login', credentials);
    return response.data;
  },

  logout: async () => {
    const response = await client.post<ApiResponse>('/auth/logout');
    return response.data;
  },

  getCurrentUser: async () => {
    const response = await client.get<ApiResponse<User>>('/auth/me');
    return response.data;
  },

  refreshToken: async (refreshToken: string) => {
    const response = await client.post<ApiResponse<{ accessToken: string }>>('/auth/refresh', {
      refreshToken,
    });
    return response.data;
  },
};
