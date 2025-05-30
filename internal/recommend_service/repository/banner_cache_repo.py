from typing import List, Dict, Optional
from model.banner import Banner

import pickle
from loguru import logger
import redis


class BannerCacheRepository:
    def __init__(self, host: str, port: int, password: Optional[str] = None):
        self.client = redis.Redis(
            host=host, port=port, password=password, decode_responses=False
        )
        self.ttl_seconds = 3 * 60  # 3 минуты в секундах

    def _cache_key(self, banner_id: int) -> str:
        return f"banner:{banner_id}"

    def get_banner(self, banner_id: int) -> Optional[Banner]:
        data = self.client.get(self._cache_key(banner_id))
        if data:
            try:
                return pickle.loads(data)
            except Exception as e:
                logger.warning(f"Failed to load banner {banner_id} from cache: {e}")
                return None
        return None

    def get_banners(self, banner_ids: List[int]) -> Dict[int, Optional[Banner]]:
        banners = {}
        pipeline = self.client.pipeline()
        keys = [self._cache_key(bid) for bid in banner_ids]

        pipeline.mget(keys)
        results = pipeline.execute()

        for bid, data in zip(banner_ids, results):
            if data:
                try:
                    banners[bid] = pickle.loads(data)
                except Exception as e:
                    logger.warning(f"Error loading banner {bid}: {e}")

        logger.debug(f"Found {len(banners)} banners in cache")
        return banners

    def set_banner(self, banner: Banner):
        key = self._cache_key(banner.id)
        self.client.setex(key, self.ttl_seconds, pickle.dumps(banner))
        logger.debug(f"Saved banner {banner.id} to cache with TTL={self.ttl_seconds}s")

    def set_banners(self, banners: list[Banner]):
        pipeline = self.client.pipeline()
        for banner in banners:
            key = self._cache_key(banner.id)
            pipeline.setex(key, self.ttl_seconds, pickle.dumps(banner))

        pipeline.execute()
        logger.debug(
            f"Saved {len(banners)} banners to cache with TTL={self.ttl_seconds}s"
        )

    def delete_banner(self, banner_id: int):
        self.client.delete(self._cache_key(banner_id))
        logger.debug(f"Deleted banner {banner_id} from cache")

    def clear_cache(self):
        self.client.flushdb()
        logger.debug("Redis cache cleared")
