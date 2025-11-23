import apiClient from './client';

export interface SetupStatus {
  setupRequired: boolean;
  adminExists: boolean;
}

export interface InitialSetupRequest {
  username: string;
  email: string;
  password: string;
  fullName: string;
}

export interface InitialSetupResponse {
  message: string;
  username: string;
  email: string;
}

export const setupApi = {
  async getStatus(): Promise<SetupStatus> {
    const response = await apiClient.get<SetupStatus>('/setup/status');
    return response.data!;
  },

  async initialize(data: InitialSetupRequest): Promise<InitialSetupResponse> {
    const response = await apiClient.post<InitialSetupResponse>('/setup/initialize', data);
    return response.data!;
  },
};
