from concurrent.futures import ThreadPoolExecutor
import grpc
from loguru import logger

import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc
import pkg.proto.banner.banner_pb2 as banner_pb2


class GrpcRecommendationServer(recommend_pb2_grpc.RecommendServiceServicer):
    def __init__(self, recommendation_service):
        self.service = recommendation_service

    def GetBannerByMetaData(self, request, context):
        logger.debug(
            f"Got request platform_id={request.platform_id}, slot_name={request.slot_name}"
        )

        banner_ids = list(request.banner_id)
        best_id = 1
        try:
            best_id = self.service.recommend_banner(banner_ids)
        except TimeoutError as e:
            logger.error(f"Timeout error during recommendation: {e}")
            context.set_details("Recommendation timed out")
            context.set_code(grpc.StatusCode.DEADLINE_EXCEEDED)
            return banner_pb2.Banner()
        except ValueError as e:
            logger.warning(f"Recommendation error: {e}")
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return banner_pb2.Banner()

        return banner_pb2.Banner(
            title="Рекламный баннер",
            content="Купи слона!",
            description="Это тестовый баннер для рекомендаций",
            link="https://example.com",
            ownerID="owner_12345",
            max_price="100.50",
            id=best_id,
        )
