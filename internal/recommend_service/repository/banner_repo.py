import psycopg2
from typing import List, Dict
from model.banner import Banner

from loguru import logger


class BannerRepository:
    def __init__(self, dsn: str):
        self.dsn = dsn

    def get_banners_by_ids(self, banner_ids: List[int]) -> Dict[int, Banner]:
        banners = {}

        with psycopg2.connect(self.dsn) as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT id, title, description, link, max_price
                    FROM banner
                    WHERE id = ANY(%s) AND NOT deleted AND status = 1
                    """,
                    (banner_ids,),
                )
                for row in cur.fetchall():
                    bid, title, desc, link, price = row
                    banners[bid] = Banner(
                        id=bid,
                        title=title,
                        description=desc,
                        link=link,
                        max_price=float(price),
                    )

        logger.debug(f"Loaded {len(banners)} banners from PostgreSQL")
        return banners
