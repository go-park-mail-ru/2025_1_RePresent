import logging
import signal
import sys
import os

PROJECT_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
sys.path.append(PROJECT_ROOT)

from internal.recommend_service.recommendApp import serve


def handle_sigterm(*_):
    logging.info("Recommend Server stopped")
    os._exit(0)


if __name__ == "__main__":
    logging.info("gRPC Recommend Server Starting on 50055 Port...")

    signal.signal(signal.SIGINT, handle_sigterm)
    signal.signal(signal.SIGTERM, handle_sigterm)

    try:
        serve()  # Однажды тут будет конфиг
    except KeyboardInterrupt:
        handle_sigterm()
