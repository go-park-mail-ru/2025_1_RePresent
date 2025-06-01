from typing import List, Dict, Optional
from model.banner import Banner, ProtoBanner
from db.connection import PostgresConnectionPool

from loguru import logger


class BannerRepository:
    def __init__(self, connection_pool: PostgresConnectionPool):
        self.connection_pool = connection_pool

    def get_banners_by_ids(self, banner_ids: List[int]) -> Dict[int, Banner]:
        banners = {}

        conn = self.connection_pool.get_connection()
        try:
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
                conn.commit()
        except Exception as e:
            conn.rollback()
            logger.error(f"Ошибка при загрузке баннеров: {e}")
        finally:
            self.connection_pool.put_connection(conn)

        logger.debug(f"Loaded {len(banners)} banners from PostgreSQL")
        return banners

    def get_proto_banner_by_id(self, banner_id: int) -> Optional[ProtoBanner]:
        """Получает один баннер по ID"""
        conn = self.connection_pool.get_connection()
        try:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT id, title, description, content, link, owner_id, max_price::text
                    FROM banner
                    WHERE id = %s AND NOT deleted AND status = 1
                    """,
                    (banner_id,),
                )
                row = cur.fetchone()

                if not row:
                    logger.warning(f"Баннер {banner_id} не найден")
                    return None

                bid, title, desc, content, link, owner_id, price = row
                return ProtoBanner(
                    id=bid,
                    title=title or "",
                    description=desc or "",
                    content=content or "",
                    link=link or "",
                    owner_id=str(owner_id),
                    max_price=price or "0.00",
                )
        except Exception as e:
            logger.error(f"Ошибка при получении баннера: {e}")
            return None
        finally:
            self.connection_pool.put_connection(conn)
