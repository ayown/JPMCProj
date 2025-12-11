"""
API response schemas
"""
from typing import Optional, Any
from pydantic import BaseModel


class HealthResponse(BaseModel):
    """Health check response"""
    status: str
    version: str
    models_loaded: bool


class ErrorResponse(BaseModel):
    """Error response"""
    error: str
    message: str
    detail: Optional[Any] = None

