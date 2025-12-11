import { apiService } from './api';
import { LoginRequest, RegisterRequest, TokenPair, User } from '@/types/auth';
import { ApiResponse } from '@/types/api';
import { STORAGE_KEYS } from '@/utils/constants';

export const authService = {
  async register(data: RegisterRequest): Promise<User> {
    const response = await apiService.post<ApiResponse<User>>('/auth/register', data);
    return response.data.data;
  },

  async login(data: LoginRequest): Promise<{ user: User; tokens: TokenPair }> {
    const response = await apiService.post<ApiResponse<TokenPair>>('/auth/login', data);
    const tokens = response.data.data;
    
    // Store tokens
    localStorage.setItem(STORAGE_KEYS.TOKEN, tokens.access_token);
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, tokens.refresh_token);
    
    // Get user profile
    const user = await this.getProfile();
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user));
    
    return { user, tokens };
  },

  async getProfile(): Promise<User> {
    const response = await apiService.get<ApiResponse<User>>('/profile');
    return response.data.data;
  },

  async refreshToken(refreshToken: string): Promise<TokenPair> {
    const response = await apiService.post<ApiResponse<TokenPair>>('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return response.data.data;
  },

  logout(): void {
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER);
  },

  getStoredUser(): User | null {
    const userStr = localStorage.getItem(STORAGE_KEYS.USER);
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch {
        return null;
      }
    }
    return null;
  },

  getStoredToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
  },

  isAuthenticated(): boolean {
    return !!this.getStoredToken();
  },
};

