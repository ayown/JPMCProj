"""
LSTM + Attention model for sequence pattern detection
"""
import os
from typing import Tuple
from app.config import settings
from app.utils.logger import log


class LSTMAttentionModel:
    """
    BiLSTM + Multi-head Attention for sequence pattern detection
    
    Captures temporal patterns in message structure
    """
    
    def __init__(self):
        self.model_path = settings.lstm_model_path
        self.loaded = False
        
        if os.path.exists(self.model_path):
            log.info(f"Loading LSTM model from {self.model_path}")
            # TODO: Load actual model
            # import tensorflow as tf
            # self.model = tf.keras.models.load_model(self.model_path)
            # self.loaded = True
            log.warning("LSTM model path exists but loading is not implemented (stub)")
        else:
            log.warning(f"LSTM model not found at {self.model_path}. Using fallback.")
    
    def predict(self, text: str) -> Tuple[float, float]:
        """
        Predict fraud probability based on sequence patterns
        
        Args:
            text: Message text
            
        Returns:
            (fraud_score, confidence)
        """
        if not self.loaded:
            return self._fallback_predict(text)
        
        # TODO: Implement actual prediction
        # sequences = self.tokenizer.texts_to_sequences([text])
        # padded = pad_sequences(sequences, maxlen=self.max_length)
        # prediction = self.model.predict(padded)
        # fraud_score = prediction[0][1]
        # confidence = max(prediction[0])
        
        return self._fallback_predict(text)
    
    def _fallback_predict(self, text: str) -> Tuple[float, float]:
        """Fallback prediction based on text patterns"""
        # Analyze text structure
        words = text.split()
        
        score = 0.0
        
        # Short messages with links are suspicious
        if len(words) < 20 and 'http' in text.lower():
            score += 0.3
        
        # Messages with many special characters
        special_chars = sum(1 for c in text if not c.isalnum() and not c.isspace())
        if special_chars / len(text) > 0.15:
            score += 0.2
        
        # Repetitive patterns (e.g., "!!!", "...")
        if '!!!' in text or '...' in text:
            score += 0.1
        
        score = min(score, 1.0)
        confidence = 0.6
        
        return score, confidence
    
    def is_loaded(self) -> bool:
        """Check if model is loaded"""
        return self.loaded


# Global instance
lstm_model = LSTMAttentionModel()

