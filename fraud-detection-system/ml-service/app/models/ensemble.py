"""
Ensemble model for fraud detection
Combines multiple models for better accuracy
"""
import os
import time
from typing import Dict, Any, Tuple
import numpy as np
from app.config import settings
from app.utils.logger import log


class FraudDetectionEnsemble:
    """
    Ensemble model that combines predictions from multiple models
    For MVP, we'll use a simple rule-based approach with feature analysis
    """
    
    def __init__(self):
        self.model_version = "1.0.0-mvp"
        self.fraud_threshold = 0.5
        self.models_loaded = False
        
        # In a full implementation, we would load actual models here
        # For MVP, we'll use rule-based detection
        log.info("Initializing Fraud Detection Ensemble (MVP mode)")
        self.models_loaded = True
    
    def predict(self, content: str, sender_header: str, features: Dict[str, Any]) -> Tuple[bool, float, str, float, Dict[str, Any], str]:
        """
        Predict if a message is fraudulent
        
        Returns:
            - is_fraud: bool
            - fraud_score: float (0-1)
            - fraud_type: str
            - confidence: float (0-1)
            - model_predictions: dict
            - explanation: str
        """
        start_time = time.time()
        
        # Extract features
        has_links = features.get("has_links", False)
        link_count = features.get("link_count", 0)
        has_urgent_words = features.get("has_urgent_words", False)
        urgent_word_count = features.get("urgent_word_count", 0)
        has_kyc_keywords = features.get("has_kyc_keywords", False)
        has_bank_names = features.get("has_bank_names", False)
        special_char_ratio = features.get("special_char_ratio", 0.0)
        capital_ratio = features.get("capital_ratio", 0.0)
        has_phone_number = features.get("has_phone_number", False)
        
        # Rule-based fraud detection (MVP approach)
        fraud_score = 0.0
        fraud_indicators = []
        
        # Check for suspicious links
        if has_links:
            fraud_score += 0.25
            fraud_indicators.append("Contains suspicious links")
        
        # Check for urgency manipulation
        if has_urgent_words:
            fraud_score += 0.15 * min(urgent_word_count / 3, 1.0)
            fraud_indicators.append(f"Uses urgent language ({urgent_word_count} urgent words)")
        
        # Check for KYC/regulatory exploitation
        if has_kyc_keywords:
            fraud_score += 0.20
            fraud_indicators.append("References KYC/regulatory requirements")
        
        # Check for excessive special characters (obfuscation)
        if special_char_ratio > 0.15:
            fraud_score += 0.10
            fraud_indicators.append("Excessive special characters")
        
        # Check for excessive capitals (shouting/urgency)
        if capital_ratio > 0.5:
            fraud_score += 0.10
            fraud_indicators.append("Excessive capital letters")
        
        # Check for phone numbers (potential callback scam)
        if has_phone_number and has_urgent_words:
            fraud_score += 0.15
            fraud_indicators.append("Contains phone number with urgent language")
        
        # Combination of bank names + urgency + links is highly suspicious
        if has_bank_names and has_urgent_words and has_links:
            fraud_score += 0.20
            fraud_indicators.append("Combines bank reference, urgency, and links")
        
        # Cap fraud score at 1.0
        fraud_score = min(fraud_score, 1.0)
        
        # Determine if it's fraud
        is_fraud = fraud_score >= self.fraud_threshold
        
        # Determine fraud type
        fraud_type = self._determine_fraud_type(features, fraud_score)
        
        # Calculate confidence
        confidence = abs(fraud_score - 0.5) * 2  # Distance from decision boundary
        confidence = max(0.6, min(confidence, 0.95))  # Clamp between 0.6 and 0.95
        
        # Generate explanation
        if is_fraud:
            explanation = f"Message flagged as fraudulent (score: {fraud_score:.2f}). Indicators: {', '.join(fraud_indicators)}"
        else:
            explanation = f"Message appears legitimate (score: {fraud_score:.2f})"
        
        # Model predictions (simulating ensemble)
        model_predictions = {
            "rule_based_model": {
                "fraud_score": fraud_score,
                "is_fraud": is_fraud,
                "indicators": fraud_indicators,
            },
            "feature_analysis": {
                "has_links": has_links,
                "has_urgent_words": has_urgent_words,
                "has_kyc_keywords": has_kyc_keywords,
                "special_char_ratio": special_char_ratio,
            }
        }
        
        inference_time = int((time.time() - start_time) * 1000)
        log.info(f"Prediction completed in {inference_time}ms: is_fraud={is_fraud}, score={fraud_score:.3f}")
        
        return is_fraud, fraud_score, fraud_type, confidence, model_predictions, explanation
    
    def _determine_fraud_type(self, features: Dict[str, Any], fraud_score: float) -> str:
        """Determine the type of fraud"""
        if fraud_score < 0.5:
            return "none"
        
        has_kyc = features.get("has_kyc_keywords", False)
        has_links = features.get("has_links", False)
        has_urgent = features.get("has_urgent_words", False)
        has_phone = features.get("has_phone_number", False)
        
        if has_kyc and has_urgent:
            return "kyc_fraud"
        elif has_links and has_urgent:
            return "phishing"
        elif has_phone and has_urgent:
            return "vishing"
        elif has_urgent:
            return "urgency_scam"
        else:
            return "generic_fraud"
    
    def is_loaded(self) -> bool:
        """Check if models are loaded"""
        return self.models_loaded


# Global ensemble instance
ensemble_model = FraudDetectionEnsemble()

