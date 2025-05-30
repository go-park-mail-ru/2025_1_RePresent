import os
from dataclasses import dataclass


@dataclass
class Config:
    db_host: str
    db_port: int
    db_user: str
    db_password: str
    db_name: str
    db_sslmode: str
    redis_host: str
    redis_port: int

    def __post_init__(self):
        self.dsn = (
            f"dbname='{self.db_name}' "
            f"user='{self.db_user}' "
            f"password='{self.db_password}' "
            f"host='{self.db_host}' "
            f"port='{self.db_port}' "
            f"sslmode='{self.db_sslmode}'"
        )


def load_config():
    return Config(
        db_host=os.getenv("DB_HOST", "localhost"),
        db_port=int(os.getenv("DB_PORT", "5432")),
        db_user=os.getenv("DB_USER", "postgres"),
        db_password=os.getenv("DB_PASSWORD", "123456"),
        db_name=os.getenv("DB_NAME", "test_db"),
        db_sslmode=os.getenv("DB_SSLMODE", "disable"),
        redis_host=os.getenv("REDIS_HOST", "localhost"),
        redis_port=int(os.getenv("REDIS_PORT", "6379")),
    )
