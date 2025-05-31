from typing import List, Dict, Optional
from model.banner import Banner
import pickle
from loguru import logger
import redis


class BannerCacheRepository:
    def __init__(self, host: str, port: int, password: Optional[str] = None):
        self.client = redis.Redis(
            host=host,
            port=port,
            password=None,
            decode_responses=False,
        )
        self.ttl_seconds = 3 * 60  # 3 минуты в секундах
        if self.client.ping():
            logger.info("Redis connection established for banners")

    def _cache_key(self, banner_id: int) -> str:
        return f"banner:{banner_id}"

    def get_banner(self, banner_id: int) -> Optional[Banner]:
        data = self.client.get(self._cache_key(banner_id))
        if data:
            try:
                banner_dict = pickle.loads(data)
                if isinstance(banner_dict, dict):
                    return Banner(**banner_dict)
                logger.warning(f"Unexpected cached data type for banner {banner_id}")
            except Exception as e:
                logger.warning(f"Failed to load banner {banner_id} from cache: {e}")
        return None

    def get_banners(self, banner_ids: List[int]) -> Dict[int, Optional[Banner]]:
        banners = {}
        keys = [self._cache_key(bid) for bid in banner_ids]

        try:
            results = self.client.mget(keys)
            for bid, data in zip(banner_ids, results):
                if data:
                    try:
                        banner_dict = pickle.loads(data)
                        if isinstance(banner_dict, dict):
                            banners[bid] = Banner(**banner_dict)
                    except Exception as e:
                        logger.warning(f"Error loading banner {bid}: {e}")
        except Exception as e:
            logger.error(f"Error getting banners from cache: {e}")

        logger.debug(f"Found {len(banners)} banners in cache")
        return banners

    def set_banner(self, banner: Banner):
        try:
            key = self._cache_key(banner.id)
            serialized = pickle.dumps(banner.to_dict())
            self.client.setex(key, self.ttl_seconds, serialized)
            logger.debug(
                f"Saved banner {banner.id} to cache with TTL={self.ttl_seconds}s"
            )
        except Exception as e:
            logger.error(f"Failed to save banner {banner.id} to cache: {e}")

    def set_banners(self, banners: List[Banner]):
        try:
            pipeline = self.client.pipeline()
            for banner in banners:
                key = self._cache_key(banner.id)
                serialized = pickle.dumps(banner.to_dict())
                pipeline.setex(key, self.ttl_seconds, serialized)
            pipeline.execute()
            logger.debug(f"Saved {len(banners)} banners to cache")
        except Exception as e:
            logger.error(f"Failed to save banners batch to cache: {e}")
