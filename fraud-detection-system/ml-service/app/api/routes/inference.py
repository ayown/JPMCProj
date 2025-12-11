"""
Inference endpoints
"""
from fastapi import APIRouter, HTTPException
from app.schemas.inference import InferenceRequest, InferenceResponse
from app.inference.predictor import predictor
from app.utils.logger import log

router = APIRouter()


@router.post("/api/v1/predict", response_model=InferenceResponse)
async def predict_fraud(request: InferenceRequest):
    """
    Predict if a message is fraudulent
    
    Args:
        request: Inference request with message content and features
    
    Returns:
        InferenceResponse with fraud prediction and details
    """
    try:
        log.info(f"Received prediction request for sender: {request.sender_header}")
        result = await predictor.predict(request)
        log.info(f"Prediction completed: is_fraud={result.is_fraud}, score={result.fraud_score:.3f}")
        return result
    except Exception as e:
        log.error(f"Prediction error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Prediction failed: {str(e)}")

