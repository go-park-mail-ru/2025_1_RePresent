from typing import List, Dict, Optional
from model.banner import Banner, ProtoBanner
from model.user import User
from repository.user_repo import UserRepository
from repository.banner_repo import BannerRepository
from repository.banner_cache_repo import BannerCacheRepository

from loguru import logger


class DataPrepareService:
    def __init__(
        self,
        user_repo: UserRepository,
        banner_repo: BannerRepository,
        banner_cache_repo: BannerCacheRepository,
    ):
        self.user_repo = user_repo
        self.banner_repo = banner_repo
        self.banner_cache_repo = banner_cache_repo

    def get_platform_by_id(self, user_id: int) -> Optional[User]:
        user = self.user_repo.get_user_by_id(user_id)

        if not user or user.role != 2:
            logger.warning(f"Platform with id:{user_id} not found")
            return None

        logger.debug(f"Platform found: {user.username} (ID: {user.id})")
        return user

    def get_cached_banners(self, banner_ids: List[int]) -> Dict[int, Optional[Banner]]:
        return self.banner_cache_repo.get_banners(banner_ids)

    def get_db_banners(self, banner_ids: List[int]) -> Dict[int, Banner]:
        return self.banner_repo.get_banners_by_ids(banner_ids)

    def get_banners(self, banner_ids: List[int]) -> Dict[int, Banner]:
        cached = self.get_cached_banners(banner_ids)
        missing = [bid for bid in banner_ids if bid not in cached]

        if missing:
            db_banners = self.get_db_banners(missing)
            for banner in db_banners.values():
                cached[banner.id] = banner
                self.banner_cache_repo.set_banner(banner)

        return cached

    def get_proto_banner_by_id(self, banner_id: int) -> ProtoBanner:
        banner = self.banner_repo.get_proto_banner_by_id(banner_id)
        if banner is None:
            raise ValueError(f"Banner {banner_id} not found")
        return banner
