"""
Main prediction service
"""
import time
import hashlib
from typing import Dict, Any
from app.models.ensemble import ensemble_model
from app.schemas.inference import InferenceRequest, InferenceResponse
from app.utils.cache import cache
from app.utils.logger import log


class FraudPredictor:
    """Main predictor service"""
    
    def __init__(self):
        self.model = ensemble_model
        self.cache_ttl = 3600  # 1 hour
    
    async def predict(self, request: InferenceRequest) -> InferenceResponse:
        """
        Predict if a message is fraudulent
        """
        start_time = time.time()
        
        # Generate cache key
        cache_key = self._generate_cache_key(request.content, request.sender_header)
        
        # Check cache
        cached_result = cache.get(cache_key)
        if cached_result:
            log.info(f"Cache hit for key: {cache_key}")
            return InferenceResponse(**cached_result)
        
        # Convert features to dict
        features_dict = request.features.dict()
        
        # Get prediction from ensemble
        is_fraud, fraud_score, fraud_type, confidence, model_predictions, explanation = self.model.predict(
            content=request.content,
            sender_header=request.sender_header,
            features=features_dict
        )
        
        # Calculate inference time
        inference_time_ms = int((time.time() - start_time) * 1000)
        
        # Create response
        response = InferenceResponse(
            is_fraud=is_fraud,
            fraud_score=fraud_score,
            fraud_type=fraud_type,
            confidence=confidence,
            model_predictions=model_predictions,
            explanation=explanation,
            inference_time_ms=inference_time_ms,
            model_version=self.model.model_version
        )
        
        # Cache the result
        cache.set(cache_key, response.dict(), ttl=self.cache_ttl)
        
        return response
    
    def _generate_cache_key(self, content: str, sender_header: str) -> str:
        """Generate cache key from content and sender"""
        combined = f"{content}:{sender_header}"
        return f"fraud_prediction:{hashlib.md5(combined.encode()).hexdigest()}"
    
    def health_check(self) -> bool:
        """Check if predictor is healthy"""
        return self.model.is_loaded()


# Global predictor instance
predictor = FraudPredictor()

