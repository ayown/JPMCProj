import { apiService } from './api';
import { ReportInput, Report, ReportStats } from '@/types/report';
import { ApiResponse } from '@/types/api';

export const reportsService = {
  async submitReport(data: ReportInput): Promise<Report> {
    const response = await apiService.post<ApiResponse<Report>>('/reports', data);
    return response.data.data;
  },

  async getReport(id: string): Promise<Report> {
    const response = await apiService.get<ApiResponse<Report>>(`/reports/${id}`);
    return response.data.data;
  },

  async getUserReports(limit: number = 20, offset: number = 0): Promise<Report[]> {
    const response = await apiService.get<ApiResponse<{ reports: Report[] }>>(
      `/reports?limit=${limit}&offset=${offset}`
    );
    return response.data.data.reports;
  },

  async getStats(): Promise<ReportStats> {
    const response = await apiService.get<ApiResponse<ReportStats>>('/reports/stats');
    return response.data.data;
  },
};

