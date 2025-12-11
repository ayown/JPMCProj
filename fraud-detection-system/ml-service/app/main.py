"""
Main FastAPI application
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.config import settings
from app.api.routes import inference_router, health_router
from app.utils.logger import log

# Create FastAPI app
app = FastAPI(
    title=settings.app_name,
    version=settings.app_version,
    description="ML service for banking fraud detection",
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(health_router, tags=["Health"])
app.include_router(inference_router, tags=["Inference"])

# Import feedback router
from app.api.routes import feedback_router
app.include_router(feedback_router, tags=["Feedback"])


@app.on_event("startup")
async def startup_event():
    """Startup event handler"""
    log.info(f"Starting {settings.app_name} v{settings.app_version}")
    log.info(f"Environment: {settings.environment}")
    log.info(f"Model path: {settings.model_path}")


@app.on_event("shutdown")
async def shutdown_event():
    """Shutdown event handler"""
    log.info("Shutting down ML service")


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "service": settings.app_name,
        "version": settings.app_version,
        "status": "running"
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.environment == "development"
    )

