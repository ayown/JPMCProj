"""
RoBERTa model for fraud type classification
"""
import os
from typing import Tuple, Dict
from app.config import settings
from app.utils.logger import log


class RoBERTaModel:
    """
    RoBERTa-based fraud type classifier
    
    Classifies fraud into types: KYC, phishing, vishing, urgency_scam, impersonation
    """
    
    def __init__(self):
        self.model_name = "roberta-base"
        self.model_path = settings.roberta_model_path
        self.loaded = False
        
        self.fraud_types = ['kyc_fraud', 'phishing', 'vishing', 'urgency_scam', 'impersonation']
        
        if os.path.exists(self.model_path):
            log.info(f"Loading RoBERTa model from {self.model_path}")
            # TODO: Load actual model
            log.warning("RoBERTa model path exists but loading is not implemented (stub)")
        else:
            log.warning(f"RoBERTa model not found at {self.model_path}. Using fallback.")
    
    def predict(self, text: str) -> Tuple[str, Dict[str, float], float]:
        """
        Predict fraud type
        
        Args:
            text: Message text
            
        Returns:
            (fraud_type, type_probabilities, confidence)
        """
        if not self.loaded:
            return self._fallback_predict(text)
        
        # TODO: Implement actual prediction
        return self._fallback_predict(text)
    
    def _fallback_predict(self, text: str) -> Tuple[str, Dict[str, float], float]:
        """Fallback prediction using rules"""
        text_lower = text.lower()
        scores = {fraud_type: 0.0 for fraud_type in self.fraud_types}
        
        # KYC fraud
        if 'kyc' in text_lower or 'know your customer' in text_lower:
            scores['kyc_fraud'] += 0.4
        
        # Phishing
        if 'http' in text_lower or 'click' in text_lower or 'link' in text_lower:
            scores['phishing'] += 0.3
        
        # Vishing
        if any(word in text_lower for word in ['call', 'phone', 'contact']):
            scores['vishing'] += 0.2
        
        # Urgency scam
        if any(word in text_lower for word in ['urgent', 'immediately', 'now', 'expire']):
            scores['urgency_scam'] += 0.3
        
        # Impersonation
        if any(word in text_lower for word in ['bank', 'rbi', 'official']):
            scores['impersonation'] += 0.2
        
        # Get top fraud type
        fraud_type = max(scores, key=scores.get)
        confidence = 0.65
        
        return fraud_type, scores, confidence
    
    def is_loaded(self) -> bool:
        """Check if model is loaded"""
        return self.loaded


# Global instance
roberta_model = RoBERTaModel()

