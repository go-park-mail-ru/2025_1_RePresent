from typing import Optional


class Banner:
    __slots__ = ("id", "title", "description", "link", "max_price")

    def __init__(
        self,
        id: int,
        title: str,
        description: str,
        link: Optional[str],
        max_price: float,
    ):
        self.id = id
        self.title = title
        self.description = description
        self.link = link
        self.max_price = max_price


class ProtoBanner:
    __slots__ = (
        "id",
        "title",
        "description",
        "content",
        "link",
        "owner_id",
        "max_price",
    )

    def __init__(
        self,
        id: int,
        title: str,
        description: str,
        content: str,
        link: str,
        owner_id: str,
        max_price: str,
    ):
        self.id = id
        self.title = title
        self.description = description
        self.content = content
        self.link = link
        self.owner_id = owner_id
        self.max_price = max_price
