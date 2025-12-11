"""
XGBoost model for metadata feature classification
"""
import os
from typing import Tuple, Dict, Any
import numpy as np
from app.config import settings
from app.utils.logger import log


class XGBoostModel:
    """
    XGBoost classifier for structured metadata features
    
    Uses non-textual features like sender info, URLs, timing, etc.
    """
    
    def __init__(self):
        self.model_path = settings.xgboost_model_path
        self.loaded = False
        
        if os.path.exists(os.path.join(self.model_path, 'model.pkl')):
            log.info(f"Loading XGBoost model from {self.model_path}")
            # TODO: Load actual model
            # import pickle
            # with open(os.path.join(self.model_path, 'model.pkl'), 'rb') as f:
            #     self.model = pickle.load(f)
            # self.loaded = True
            log.warning("XGBoost model path exists but loading is not implemented (stub)")
        else:
            log.warning(f"XGBoost model not found at {self.model_path}. Using fallback.")
    
    def predict(self, features: Dict[str, Any]) -> Tuple[float, float]:
        """
        Predict fraud probability based on metadata features
        
        Args:
            features: Dictionary of extracted features
            
        Returns:
            (fraud_score, confidence)
        """
        if not self.loaded:
            return self._fallback_predict(features)
        
        # TODO: Implement actual prediction
        # feature_vector = self._extract_feature_vector(features)
        # fraud_prob = self.model.predict_proba([feature_vector])[0][1]
        # confidence = max(self.model.predict_proba([feature_vector])[0])
        
        return self._fallback_predict(features)
    
    def _fallback_predict(self, features: Dict[str, Any]) -> Tuple[float, float]:
        """Fallback prediction using feature analysis"""
        score = 0.0
        
        # Link-based features
        if features.get('has_links', False):
            score += 0.2
            score += features.get('link_count', 0) * 0.1
        
        # Phone number features
        if features.get('has_phone_number', False):
            score += 0.15
        
        # Text ratio features
        special_ratio = features.get('special_char_ratio', 0)
        if special_ratio > 0.15:
            score += 0.15
        
        capital_ratio = features.get('capital_ratio', 0)
        if capital_ratio > 0.5:
            score += 0.1
        
        # Urgency features
        if features.get('has_urgent_words', False):
            score += 0.2
        
        # KYC features
        if features.get('has_kyc_keywords', False):
            score += 0.15
        
        score = min(score, 1.0)
        confidence = 0.75  # Higher confidence for structured features
        
        return score, confidence
    
    def _extract_feature_vector(self, features: Dict[str, Any]) -> np.ndarray:
        """Extract feature vector for XGBoost"""
        # TODO: Implement proper feature extraction
        feature_names = [
            'message_length', 'has_links', 'link_count', 'has_phone_number',
            'phone_number_count', 'has_urgent_words', 'urgent_word_count',
            'special_char_ratio', 'capital_ratio', 'number_ratio',
            'has_kyc_keywords', 'has_bank_names'
        ]
        
        vector = []
        for name in feature_names:
            value = features.get(name, 0)
            if isinstance(value, bool):
                value = 1 if value else 0
            vector.append(float(value))
        
        return np.array(vector)
    
    def is_loaded(self) -> bool:
        """Check if model is loaded"""
        return self.loaded


# Global instance
xgboost_model = XGBoostModel()

