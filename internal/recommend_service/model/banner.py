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
