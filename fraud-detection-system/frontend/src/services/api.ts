import axios, { AxiosInstance, AxiosError } from 'axios';
import { API_URL, STORAGE_KEYS } from '@/utils/constants';
import toast from 'react-hot-toast';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor
    this.api.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.api.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Token expired, try to refresh
          const refreshToken = localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
          if (refreshToken) {
            try {
              const response = await this.api.post('/auth/refresh', {
                refresh_token: refreshToken,
              });
              const { access_token, refresh_token: newRefreshToken } = response.data.data;
              localStorage.setItem(STORAGE_KEYS.TOKEN, access_token);
              localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, newRefreshToken);
              
              // Retry original request
              if (error.config) {
                error.config.headers.Authorization = `Bearer ${access_token}`;
                return this.api.request(error.config);
              }
            } catch (refreshError) {
              // Refresh failed, logout
              this.logout();
              toast.error('Session expired. Please login again.');
            }
          } else {
            this.logout();
          }
        }
        return Promise.reject(error);
      }
    );
  }

  private logout() {
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER);
    window.location.href = '/login';
  }

  get<T>(url: string, config?: any) {
    return this.api.get<T>(url, config);
  }

  post<T>(url: string, data?: any, config?: any) {
    return this.api.post<T>(url, data, config);
  }

  put<T>(url: string, data?: any, config?: any) {
    return this.api.put<T>(url, data, config);
  }

  delete<T>(url: string, config?: any) {
    return this.api.delete<T>(url, config);
  }
}

export const apiService = new ApiService();

