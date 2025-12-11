export interface ReportInput {
  message_id?: string;
  verification_id?: string;
  report_type: 'FRAUD' | 'FALSE_POSITIVE' | 'FEEDBACK';
  content: string;
  sender_header: string;
  description: string;
}

export interface Report {
  id: string;
  report_type: string;
  content: string;
  sender_header: string;
  description: string;
  status: 'PENDING' | 'REVIEWED' | 'RESOLVED' | 'DISMISSED';
  priority: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  reviewed_at?: string;
  review_notes?: string;
  created_at: string;
}

export interface ReportStats {
  total_reports: number;
  pending_reports: number;
  resolved_reports: number;
  by_type: Record<string, number>;
  by_priority: Record<string, number>;
  last_24_hours: number;
  last_7_days: number;
}

