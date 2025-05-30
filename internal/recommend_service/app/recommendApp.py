import grpc
from concurrent.futures import ThreadPoolExecutor

import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc
import pkg.proto.banner.banner_pb2 as banner_pb2
from controller.grpc.recommend_handler import GrpcRecommendationServer
from service.recommendation_service import RecommendationService
from repository.user_repo import UserRepository
from repository.banner_repo import BannerRepository
from repository.banner_cache_repo import BannerCacheRepository
from config import load_config

from loguru import logger


def serve(config=None):
    if config is None:
        config = load_config()

    dsn = (
        f"dbname='{config.db_name}' "
        f"user='{config.db_user}' "
        f"password='{config.db_password}' "
        f"host='{config.db_host}' "
        f"port='{config.db_port}' "
        f"sslmode='{config.db_sslmode}'"
    )

    user_repo = UserRepository(dsn)
    banner_repo = BannerRepository(dsn)
    banner_cache = BannerCacheRepository(config.redis_host, config.redis_port)

    recommendation_service = RecommendationService(user_repo, banner_repo, banner_cache)

    server = grpc.server(ThreadPoolExecutor(max_workers=10))

    recommend_pb2_grpc.add_RecommendServiceServicer_to_server(
        GrpcRecommendationServer(recommendation_service), server
    )

    server.add_insecure_port("[::]:50055")
    logger.info("gRPC Recommend Server Started on ReTargetApiRecommend:50055")
    server.start()

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Stopping Recommend Server...")
        server.stop(grace=5)
        logger.info("Recommend Server stopped")
