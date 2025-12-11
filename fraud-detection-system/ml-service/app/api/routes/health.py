"""
Health check endpoints
"""
from fastapi import APIRouter
from app.schemas.response import HealthResponse
from app.config import settings
from app.inference.predictor import predictor
from app.utils.cache import cache

router = APIRouter()


@router.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        version=settings.app_version,
        models_loaded=predictor.health_check()
    )


@router.get("/ready")
async def readiness_check():
    """Readiness check endpoint"""
    models_ready = predictor.health_check()
    cache_ready = cache.health_check()
    
    if models_ready:
        return {
            "status": "ready",
            "models": "loaded",
            "cache": "connected" if cache_ready else "disconnected"
        }
    else:
        return {
            "status": "not ready",
            "models": "not loaded",
            "cache": "connected" if cache_ready else "disconnected"
        }, 503

