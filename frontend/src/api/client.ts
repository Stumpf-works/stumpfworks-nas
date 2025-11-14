import axios, { AxiosError, AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: number;
    message: string;
  };
}

// Extended request config to track retries
interface RetryConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

// Create axios instance
const client: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor (add auth token)
client.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken');
    console.log('[API Client] Request interceptor - Token exists:', !!token);
    console.log('[API Client] Request URL:', config.url);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
      console.log('[API Client] Authorization header added');
    } else {
      console.warn('[API Client] NO TOKEN FOUND in localStorage!');
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor (handle errors)
client.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    return response;
  },
  async (error: AxiosError<ApiResponse>) => {
    const originalRequest = error.config as RetryConfig | undefined;

    // If 401 and we have a refresh token, try to refresh
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem('refreshToken');
      if (refreshToken) {
        try {
          const response = await axios.post<ApiResponse<{ accessToken: string }>>(
            '/api/v1/auth/refresh',
            { refreshToken }
          );

          if (response.data.success && response.data.data?.accessToken) {
            localStorage.setItem('accessToken', response.data.data.accessToken);

            // Retry original request with new token
            if (originalRequest) {
              originalRequest.headers.Authorization = `Bearer ${response.data.data.accessToken}`;
              return client.request(originalRequest);
            }
          }
        } catch (refreshError) {
          // Refresh failed, logout
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          window.location.href = '/login';
          return Promise.reject(refreshError);
        }
      } else {
        // No refresh token, redirect to login
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

export default client;

// Helper function to handle API errors
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const apiError = error as AxiosError<ApiResponse>;
    return apiError.response?.data?.error?.message || error.message || 'An error occurred';
  }
  if (error instanceof Error) {
    return error.message;
  }
  return 'An unknown error occurred';
}
