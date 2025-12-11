import { apiService } from './api';
import { VerificationRequest, VerificationResponse, VerificationStats } from '@/types/verification';
import { ApiResponse } from '@/types/api';

export const verificationService = {
  async verifyMessage(data: VerificationRequest): Promise<VerificationResponse> {
    const response = await apiService.post<ApiResponse<VerificationResponse>>('/verify', data);
    return response.data.data;
  },

  async getVerification(id: string): Promise<VerificationResponse> {
    const response = await apiService.get<ApiResponse<VerificationResponse>>(`/verify/${id}`);
    return response.data.data;
  },

  async getHistory(limit: number = 20, offset: number = 0): Promise<VerificationResponse[]> {
    const response = await apiService.get<ApiResponse<{ verifications: VerificationResponse[] }>>(
      `/verify/history?limit=${limit}&offset=${offset}`
    );
    return response.data.data.verifications;
  },

  async getStats(): Promise<VerificationStats> {
    const response = await apiService.get<ApiResponse<VerificationStats>>('/verify/stats');
    return response.data.data;
  },
};

