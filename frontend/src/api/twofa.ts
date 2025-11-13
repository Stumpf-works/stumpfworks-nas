import { apiClient } from './client';

export interface TwoFAStatus {
  enabled: boolean;
  backupCodesRemaining: number;
}

export interface TwoFASetupResponse {
  secret: string;
  qrCodeUrl: string;
  backupCodes: string[];
}

export interface TwoFAVerifyRequest {
  code: string;
}

export interface TwoFALoginRequest {
  userId: number;
  code: string;
  isBackupCode: boolean;
}

export const twofaApi = {
  /**
   * Get 2FA status for the current user
   */
  getStatus: async (): Promise<TwoFAStatus> => {
    const response = await apiClient.get('/api/v1/2fa/status');
    return response.data.data;
  },

  /**
   * Setup 2FA (returns QR code and backup codes)
   */
  setup: async (): Promise<TwoFASetupResponse> => {
    const response = await apiClient.post('/api/v1/2fa/setup');
    return response.data.data;
  },

  /**
   * Enable 2FA after setup (requires verification code)
   */
  enable: async (code: string): Promise<void> => {
    await apiClient.post('/api/v1/2fa/enable', { code });
  },

  /**
   * Disable 2FA (requires verification code)
   */
  disable: async (code: string): Promise<void> => {
    await apiClient.post('/api/v1/2fa/disable', { code });
  },

  /**
   * Regenerate backup codes (requires verification code)
   */
  regenerateBackupCodes: async (code: string): Promise<string[]> => {
    const response = await apiClient.post('/api/v1/2fa/backup-codes/regenerate', {
      code,
    });
    return response.data.data.backupCodes;
  },

  /**
   * Complete login with 2FA code
   */
  loginWith2FA: async (req: TwoFALoginRequest): Promise<any> => {
    const response = await apiClient.post('/api/v1/auth/login/2fa', req);
    return response.data.data;
  },
};
