import grpc
from concurrent.futures import ThreadPoolExecutor

import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc
import pkg.proto.banner.banner_pb2 as banner_pb2
from db.connection import PostgresConnectionPool
from controller.grpc.recommend_handler import GrpcRecommendationServer
from service.recommendation_service import RecommendationService
from service.deta_prepare_service import DataPrepareService

from repository.user_repo import UserRepository
from repository.banner_repo import BannerRepository
from repository.banner_cache_repo import BannerCacheRepository
from repository.embedding_cache_repo import EmbeddingCacheRepository
from config import load_config

from loguru import logger


def serve(config=None):
    if config is None:
        config = load_config()

    connection_pool = PostgresConnectionPool(dsn=config.dsn, minconn=1, maxconn=5)

    user_repo = UserRepository(connection_pool)
    banner_repo = BannerRepository(connection_pool)
    banner_cache = BannerCacheRepository(
        config.redis_host, config.redis_port, password=config.redis_password
    )
    embedding_cache = EmbeddingCacheRepository(
        config.redis_host, config.redis_port, password=config.redis_password
    )

    data_prepare_service = DataPrepareService(user_repo, banner_repo, banner_cache)
    recommendation_service = RecommendationService(embedding_cache)

    server = grpc.server(ThreadPoolExecutor(max_workers=10))

    recommend_pb2_grpc.add_RecommendServiceServicer_to_server(
        GrpcRecommendationServer(recommendation_service, data_prepare_service), server
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
