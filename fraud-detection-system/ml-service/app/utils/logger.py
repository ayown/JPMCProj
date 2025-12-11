"""
Logging utilities
"""
import sys
from loguru import logger
from app.config import settings


def setup_logger():
    """Setup logger configuration"""
    logger.remove()
    logger.add(
        sys.stdout,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
        level=settings.log_level,
        colorize=True,
    )
    return logger


# Initialize logger
log = setup_logger()

