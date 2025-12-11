"""
Inference request and response schemas
"""
from typing import List, Optional, Dict, Any
from pydantic import BaseModel, Field


class MessageFeatures(BaseModel):
    """Features extracted from a message"""
    content: str
    sender_header: str
    message_length: int
    has_links: bool
    link_count: int
    extracted_urls: List[str] = []
    has_phone_number: bool
    phone_number_count: int
    has_urgent_words: bool
    urgent_word_count: int
    special_char_ratio: float
    capital_ratio: float
    number_ratio: float
    has_kyc_keywords: bool
    has_bank_names: bool


class InferenceRequest(BaseModel):
    """Request for fraud detection inference"""
    content: str = Field(..., description="Message content to analyze")
    sender_header: str = Field(..., description="Sender header/ID")
    features: MessageFeatures = Field(..., description="Extracted message features")


class InferenceResponse(BaseModel):
    """Response from fraud detection inference"""
    is_fraud: bool = Field(..., description="Whether the message is fraudulent")
    fraud_score: float = Field(..., description="Fraud probability score (0-1)")
    fraud_type: str = Field(..., description="Type of fraud detected")
    confidence: float = Field(..., description="Model confidence (0-1)")
    model_predictions: Dict[str, Any] = Field(..., description="Individual model predictions")
    explanation: str = Field(..., description="Human-readable explanation")
    inference_time_ms: int = Field(..., description="Inference time in milliseconds")
    model_version: str = Field(..., description="Model version used")

