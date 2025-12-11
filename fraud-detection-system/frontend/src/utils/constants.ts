export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
export const APP_NAME = import.meta.env.VITE_APP_NAME || 'Fraud Detection System';

export const RISK_LEVELS = {
  LOW: {
    color: 'success',
    label: 'Low Risk',
    bgColor: 'bg-success-100',
    textColor: 'text-success-800',
    borderColor: 'border-success-500',
  },
  MEDIUM: {
    color: 'warning',
    label: 'Medium Risk',
    bgColor: 'bg-yellow-100',
    textColor: 'text-yellow-800',
    borderColor: 'border-yellow-500',
  },
  HIGH: {
    color: 'danger',
    label: 'High Risk',
    bgColor: 'bg-danger-100',
    textColor: 'text-danger-800',
    borderColor: 'border-danger-500',
  },
  CRITICAL: {
    color: 'danger',
    label: 'Critical Risk',
    bgColor: 'bg-danger-200',
    textColor: 'text-danger-900',
    borderColor: 'border-danger-700',
  },
} as const;

export const FRAUD_TYPES = {
  kyc_fraud: 'KYC Fraud',
  phishing: 'Phishing',
  vishing: 'Vishing',
  urgency_scam: 'Urgency Scam',
  impersonation: 'Impersonation',
  generic_fraud: 'Generic Fraud',
  none: 'None',
} as const;

export const REPORT_TYPES = {
  FRAUD: 'Fraud Report',
  FALSE_POSITIVE: 'False Positive',
  FEEDBACK: 'Feedback',
} as const;

export const STORAGE_KEYS = {
  TOKEN: 'fraud_detection_token',
  REFRESH_TOKEN: 'fraud_detection_refresh_token',
  USER: 'fraud_detection_user',
} as const;

