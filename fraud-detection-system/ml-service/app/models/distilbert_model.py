"""
DistilBERT model for fraud classification
"""
import os
from typing import Tuple, List
import numpy as np
from app.config import settings
from app.utils.logger import log


class DistilBERTModel:
    """
    DistilBERT-based fraud classifier
    
    Note: This is a stub implementation. To use actual DistilBERT:
    1. Train the model using training/train_distilbert.py
    2. Save the model to models/distilbert/
    3. Uncomment the transformers imports below
    4. Load the model in __init__
    """
    
    def __init__(self):
        self.model_name = "distilbert-base-uncased"
        self.model_path = settings.distilbert_model_path
        self.loaded = False
        
        # TODO: Uncomment when model is trained
        # from transformers import DistilBertTokenizer, DistilBertForSequenceClassification
        # import torch
        
        # Check if model exists
        if os.path.exists(self.model_path):
            log.info(f"Loading DistilBERT model from {self.model_path}")
            # TODO: Load actual model
            # self.tokenizer = DistilBertTokenizer.from_pretrained(self.model_path)
            # self.model = DistilBertForSequenceClassification.from_pretrained(self.model_path)
            # self.model.eval()
            # self.loaded = True
            log.warning("DistilBERT model path exists but loading is not implemented (stub)")
        else:
            log.warning(f"DistilBERT model not found at {self.model_path}. Using fallback.")
    
    def predict(self, text: str) -> Tuple[float, float]:
        """
        Predict fraud probability
        
        Args:
            text: Message text
            
        Returns:
            (fraud_score, confidence)
        """
        if not self.loaded:
            # Fallback: Simple keyword-based prediction
            return self._fallback_predict(text)
        
        # TODO: Implement actual prediction
        # inputs = self.tokenizer(text, return_tensors="pt", truncation=True, max_length=512)
        # with torch.no_grad():
        #     outputs = self.model(**inputs)
        #     probs = torch.softmax(outputs.logits, dim=1)
        #     fraud_prob = probs[0][1].item()
        #     confidence = max(probs[0]).item()
        # return fraud_prob, confidence
        
        return self._fallback_predict(text)
    
    def _fallback_predict(self, text: str) -> Tuple[float, float]:
        """Fallback prediction using simple rules"""
        text_lower = text.lower()
        score = 0.0
        
        fraud_keywords = ['urgent', 'blocked', 'suspended', 'kyc', 'verify', 'click', 'link']
        for keyword in fraud_keywords:
            if keyword in text_lower:
                score += 0.15
        
        score = min(score, 1.0)
        confidence = 0.7  # Lower confidence for fallback
        
        return score, confidence
    
    def is_loaded(self) -> bool:
        """Check if model is loaded"""
        return self.loaded


# Global instance
distilbert_model = DistilBERTModel()

