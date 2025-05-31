import signal
import sys
import os

from loguru import logger

logger.remove()
logger.add(
    sys.stdout, level="INFO", format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {message}"
)
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..")))
sys.path.append(
    os.path.abspath(
        os.path.join(
            os.path.dirname(__file__), "..", "..", "internal", "recommend_service"
        )
    )
)

from internal.recommend_service.app.recommendApp import serve


def handle_sigterm(*_):
    logger.info("Recommend Server stopped")
    os._exit(0)


if __name__ == "__main__":
    logger.info("gRPC Recommend Server Starting on 50055 Port...")

    signal.signal(signal.SIGINT, handle_sigterm)
    signal.signal(signal.SIGTERM, handle_sigterm)

    try:
        serve()  # Однажды тут будет конфиг
    except KeyboardInterrupt:
        handle_sigterm()
