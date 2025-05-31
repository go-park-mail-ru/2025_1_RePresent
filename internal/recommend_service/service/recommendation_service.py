from typing import List, Dict, Tuple, Optional
import numpy as np
from loguru import logger
from sentence_transformers import SentenceTransformer
import faiss
from catboost import CatBoostRanker, Pool
from utils.timeout import with_timeout

from model.banner import Banner
from repository.embedding_cache_repo import EmbeddingCacheRepository


class RecommendationService:
    def __init__(self, emcache_repo: EmbeddingCacheRepository):
        self.model = SentenceTransformer("distiluse-base-multilingual-cased-v2")
        self.ranker = CatBoostRanker()
        self._load_ranker_model()
        self.emcache_repo = emcache_repo

    def _load_ranker_model(self):
        try:
            self.ranker.load_model("reluma.cbm")
            logger.info("ReLuma model loaded")
        except Exception as e:
            logger.warning(
                f"Failed to load ReLuma model: {e}, fallback to default scoring"
            )

    def _create_text_embedding(
        self, text: str, banner_id: Optional[int], local_cache: dict
    ) -> np.ndarray:
        if banner_id is not None and banner_id in local_cache:
            return local_cache[banner_id]

        emb = self.model.encode([text], convert_to_numpy=True)
        norm_emb = emb / np.linalg.norm(emb)

        if banner_id is not None:
            local_cache[banner_id] = norm_emb[0]

        return norm_emb[0]

    def _build_query_embedding(
        self,
        platform_title: str,
        platform_description: str,
        slot_name: str,
        local_cache: dict,
    ) -> np.ndarray:
        combined_text = (
            f"{slot_name} " * 5 + f"{platform_description} " * 4 + platform_title
        )
        return self._create_text_embedding(combined_text, None, local_cache)

    def _prepare_candidate_features(
        self,
        query_emb: np.ndarray,
        banners: Dict[int, Banner],
        candidate_ids: List[int],
        local_cache: dict,
    ) -> Tuple[List[Dict], List[int]]:
        features = []
        ids = []

        cached_embs = self.emcache_repo.get_embeddings(candidate_ids)
        new_embeddings = {}

        for bid in candidate_ids:
            banner = banners[bid]
            title = str(banner.title or "")
            description = str(banner.description or "")

            key = bid
            if key in cached_embs:
                banner_emb = cached_embs[key]
            else:
                banner_emb = self._create_text_embedding(
                    f"{title} {description}", bid, local_cache
                )
                new_embeddings[key] = banner_emb

            similarity = float(np.dot(query_emb, banner_emb))
            price_weighted_sim = similarity * 0.45 + banner.max_price * 0.55

            features.append(
                {
                    "similarity": similarity,
                    "max_price": float(banner.max_price),
                    "title_len": len(title),
                    "desc_len": len(description),
                    "price_weighted_sim": price_weighted_sim,
                }
            )
            ids.append(bid)

        if new_embeddings:
            self.emcache_repo.set_embeddings(new_embeddings)

        return features, ids

    def _build_temporary_index(
        self, banners: Dict[int, Banner], local_cache: dict
    ) -> Tuple[faiss.Index, List[int]]:
        texts = [f"{b.title or ''} {b.description or ''}" for b in banners.values()]
        embeddings = np.vstack(
            [
                self._create_text_embedding(text, bid, local_cache)
                for bid, text in zip(banners.keys(), texts)
            ]
        )
        index = faiss.IndexFlatIP(self.model.get_sentence_embedding_dimension())
        index.add(embeddings)
        return index, list(banners.keys())

    @with_timeout(seconds=1.5)
    def recommend_banner(
        self,
        slot_name: str,
        paltform_username: str,
        platform_description: str,
        banners: Dict[int, Banner],
    ) -> int:
        if not banners:
            raise ValueError("Empty banners for recommend")

        local_cache = {}

        query_emb = self._build_query_embedding(
            paltform_username, platform_description, slot_name, local_cache
        )

        temp_index, banner_ids = self._build_temporary_index(banners, local_cache)
        TOP_K = min(20, len(banners))
        D, I = temp_index.search(np.expand_dims(query_emb, axis=0), k=TOP_K)
        candidate_indices = I[0]
        candidate_ids = [banner_ids[i] for i in candidate_indices]

        feature_order = [
            "similarity",
            "max_price",
            "title_len",
            "desc_len",
            "price_weighted_sim",
        ]
        features, feature_ids = self._prepare_candidate_features(
            query_emb, banners, candidate_ids, local_cache
        )
        if not self.ranker.is_fitted():
            scores = [f["price_weighted_sim"] for f in features]
            best_idx = int(np.argmax(scores))
        else:
            feature_list = [[f[k] for k in feature_order] for f in features]
            pool = Pool(data=feature_list)
            preds = self.ranker.predict(pool)
            best_idx = int(np.argmax(preds))

        # 10 %
        tolerance = 0.1
        best_score = preds[best_idx]

        candidates_within_range = [
            idx
            for idx, score in enumerate(preds)
            if score >= best_score * (1 - tolerance)
        ]

        if not candidates_within_range:
            candidates_within_range = [best_idx]

        final_idx = np.random.choice(candidates_within_range)
        return feature_ids[final_idx]
