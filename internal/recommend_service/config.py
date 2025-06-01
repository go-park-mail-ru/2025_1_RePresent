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
    redis_password: str

    def __post_init__(self):
        self.dsn = (
            f"dbname='{self.db_name}' "
            f"user='{self.db_user}' "
            f"password='{self.db_password}' "
            f"host='{self.db_host}' "
            f"port={self.db_port} "
            f"sslmode='{self.db_sslmode}'"
        )


def load_config():
    return Config(
        db_host=os.getenv("PSQL_HOST", "ReTargetDataBase"),
        db_port=int(os.getenv("PSQL_INSIDE_PORT", "5432")),
        db_user=os.getenv("PSQL_USER", "postgres"),
        db_password=os.getenv("PSQL_PASSWORD", "123456"),
        db_name=os.getenv("PSQL_DB_NAME", "test_db"),
        db_sslmode=os.getenv("PSQL_SSLMODE", "disable"),
        redis_host=os.getenv("REDIS_HOST", "ReTargetRedis"),
        redis_port=int(os.getenv("REDIS_PORT", "6379")),
        redis_password=os.getenv("REDIS_PASSWORD", None),
    )
