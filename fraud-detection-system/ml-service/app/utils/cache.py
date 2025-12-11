"""
Redis caching utilities
"""
import json
from typing import Optional, Any
import redis
from app.config import settings
from app.utils.logger import log


class RedisCache:
    """Redis cache client"""
    
    def __init__(self):
        try:
            self.client = redis.Redis(
                host=settings.redis_host,
                port=settings.redis_port,
                db=settings.redis_db,
                password=settings.redis_password,
                decode_responses=True,
            )
            self.client.ping()
            log.info("Redis cache connected")
        except Exception as e:
            log.warning(f"Redis connection failed: {e}. Caching disabled.")
            self.client = None
    
    def get(self, key: str) -> Optional[Any]:
        """Get value from cache"""
        if not self.client:
            return None
        
        try:
            value = self.client.get(key)
            if value:
                return json.loads(value)
        except Exception as e:
            log.error(f"Cache get error: {e}")
        return None
    
    def set(self, key: str, value: Any, ttl: int = 3600):
        """Set value in cache"""
        if not self.client:
            return
        
        try:
            self.client.setex(key, ttl, json.dumps(value))
        except Exception as e:
            log.error(f"Cache set error: {e}")
    
    def delete(self, key: str):
        """Delete value from cache"""
        if not self.client:
            return
        
        try:
            self.client.delete(key)
        except Exception as e:
            log.error(f"Cache delete error: {e}")
    
    def health_check(self) -> bool:
        """Check if cache is healthy"""
        if not self.client:
            return False
        
        try:
            return self.client.ping()
        except:
            return False


# Global cache instance
cache = RedisCache()

