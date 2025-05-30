import psycopg2
from model.user import User
from typing import Dict, List, Optional

from loguru import logger


class UserRepository:
    def __init__(self, dsn: str):
        self.dsn = dsn

    def get_user_by_id(self, user_id: int) -> Optional[User]:
        with psycopg2.connect(self.dsn) as conn:
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

    def get_users_by_ids(self, user_ids: List[int]) -> Dict[int, User]:
        users = {}
        if not user_ids:
            return users

        with psycopg2.connect(self.dsn) as conn:
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
        logger.debug(f"Loaded {len(users)} users from DB")
        return users
