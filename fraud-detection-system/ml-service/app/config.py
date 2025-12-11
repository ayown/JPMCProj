"""
Configuration for ML Service
"""
import os
from typing import Optional
from pydantic import BaseSettings


class Settings(BaseSettings):
    # Application
    app_name: str = "Fraud Detection ML Service"
    app_version: str = "1.0.0"
    environment: str = os.getenv("APP_ENV", "development")
    
    # Server
    host: str = os.getenv("ML_SERVICE_HOST", "0.0.0.0")
    port: int = int(os.getenv("ML_SERVICE_PORT", "8000"))
    
    # Model paths
    model_path: str = os.getenv("MODEL_PATH", "/app/models")
    distilbert_model_path: str = os.getenv("DISTILBERT_MODEL_PATH", "/app/models/distilbert")
    roberta_model_path: str = os.getenv("ROBERTA_MODEL_PATH", "/app/models/roberta")
    lstm_model_path: str = os.getenv("LSTM_MODEL_PATH", "/app/models/lstm")
    xgboost_model_path: str = os.getenv("XGBOOST_MODEL_PATH", "/app/models/xgboost")
    ensemble_model_path: str = os.getenv("ENSEMBLE_MODEL_PATH", "/app/models/ensemble")
    
    # Redis
    redis_host: str = os.getenv("REDIS_HOST", "localhost")
    redis_port: int = int(os.getenv("REDIS_PORT", "6379"))
    redis_db: int = int(os.getenv("REDIS_DB", "0"))
    redis_password: Optional[str] = os.getenv("REDIS_PASSWORD")
    
    # Model settings
    max_length: int = 512
    batch_size: int = 32
    confidence_threshold: float = 0.5
    
    # Hugging Face
    huggingface_token: Optional[str] = os.getenv("HUGGINGFACE_TOKEN")
    
    # Logging
    log_level: str = os.getenv("LOG_LEVEL", "INFO")
    
    class Config:
        env_file = ".env"
        case_sensitive = False


settings = Settings()

