import psycopg2
from psycopg2 import pool
from loguru import logger


class PostgresConnectionPool:
    def __init__(self, dsn: str, minconn=1, maxconn=5):
        self.dsn = dsn
        self.pool = psycopg2.pool.ThreadedConnectionPool(
            minconn=minconn, maxconn=maxconn, dsn=dsn
        )
        logger.info("PostgreSQL connection pool created")

    def get_connection(self):
        return self.pool.getconn()

    def put_connection(self, conn):
        self.pool.putconn(conn)

    def close_all(self):
        self.pool.closeall()
        logger.info("All PSQL connections closed")
