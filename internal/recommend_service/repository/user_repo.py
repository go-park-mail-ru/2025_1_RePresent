from typing import List, Dict, Optional
from model.user import User
from db.connection import PostgresConnectionPool

from loguru import logger


class UserRepository:
    def __init__(self, connection_pool: PostgresConnectionPool):
        self.connection_pool = connection_pool

    def get_user_by_id(self, user_id: int) -> Optional[User]:
        conn = self.connection_pool.get_connection()
        try:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT id, username, description, role
                    FROM auth_user
                    WHERE id = %s AND NOT deleted
                    """,
                    (user_id,),
                )
                row = cur.fetchone()

                if not row:
                    logger.debug(f"User {user_id} not found")
                    return None

                uid, username, description, role = row
                return User(
                    id=uid, username=username, description=description or "", role=role
                )
        except Exception as e:
            logger.error(f"Ошибка при получении пользователя {user_id}: {e}")
            return None
        finally:
            self.connection_pool.put_connection(conn)

    def get_users_by_ids(self, user_ids: List[int]) -> Dict[int, User]:
        users = {}
        if not user_ids:
            return users

        conn = self.connection_pool.get_connection()
        try:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT id, username, description, role
                    FROM auth_user
                    WHERE id = ANY(%s)
                    """,
                    (user_ids,),
                )
                for row in cur.fetchall():
                    uid, username, description, role = row
                    users[uid] = User(
                        id=uid,
                        username=username,
                        description=description or "",
                        role=role,
                    )
                conn.commit()
        except Exception as e:
            conn.rollback()
            logger.error(f"Ошибка при загрузке пользователей: {e}")
        finally:
            self.connection_pool.put_connection(conn)

        logger.debug(f"Loaded {len(users)} users from DB")
        return users
