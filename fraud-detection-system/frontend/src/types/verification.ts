export interface VerificationRequest {
  content: string;
  sender_header: string;
  message_type?: 'SMS' | 'WhatsApp' | 'Email';
  received_at?: string;
  phone_number?: string;
}

export interface VerificationResponse {
  id: string;
  message_id: string;
  is_fraud: boolean;
  fraud_score: number;
  fraud_type?: string;
  confidence: number;
  risk_level: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  header_verified: boolean;
  rbi_compliant: boolean;
  explanation: string;
  recommendations: string[];
  model_predictions: Record<string, any>;
  processing_time_ms: number;
  verified_at: string;
}

export interface VerificationStats {
  total_verifications: number;
  fraud_detected: number;
  fraud_rate: number;
  avg_fraud_score: number;
  avg_processing_time: number;
  last_24_hours: number;
  last_7_days: number;
}

