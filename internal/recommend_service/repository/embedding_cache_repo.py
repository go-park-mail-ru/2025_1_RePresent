from typing import Dict, List, Optional
from loguru import logger
import numpy as np
import redis
import pickle


class EmbeddingCacheRepository:
    def __init__(self, host: str, port: int, password: Optional[str] = None):
        self.client = redis.Redis(
            host=host, port=port, password=None, decode_responses=False
        )
        self.ttl_seconds = 3 * 60
        if self.client.ping():
            logger.info("Redis connection established for embeddings")

    def _cache_key(self, banner_id: int) -> str:
        return f"embedding:{banner_id}"

    def get_embeddings(self, banner_ids: List[int]) -> Dict[int, np.ndarray]:
        keys = [self._cache_key(bid) for bid in banner_ids]
        results = self.client.mget(keys)
        embeddings = {}
        for bid, data in zip(banner_ids, results):
            if data:
                try:
                    emb = pickle.loads(data)
                    if isinstance(emb, list):
                        emb = np.array(emb)
                    embeddings[bid] = emb
                except Exception:
                    pass
        return embeddings

    def set_embeddings(self, embeddings: Dict[int, np.ndarray]):
        pipeline = self.client.pipeline()
        for bid, emb in embeddings.items():
            key = self._cache_key(bid)
            pipeline.setex(key, self.ttl_seconds, pickle.dumps(emb.tolist()))
        pipeline.execute()
