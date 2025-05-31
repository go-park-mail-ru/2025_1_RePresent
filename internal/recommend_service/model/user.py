from typing import Optional


class User:
    __slots__ = ("id", "username", "description", "role")

    def __init__(
        self,
        id: int,
        username: str,
        description: str,
        role: int,  # 1=advertiser, 2=platform
    ):
        self.id = id
        self.username = username
        self.description = description
        self.role = role
