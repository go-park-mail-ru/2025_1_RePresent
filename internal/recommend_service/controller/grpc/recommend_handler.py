from concurrent.futures import ThreadPoolExecutor
import grpc
from loguru import logger

import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc
import pkg.proto.banner.banner_pb2 as banner_pb2


class GrpcRecommendationServer(recommend_pb2_grpc.RecommendServiceServicer):
    def __init__(self, recommendation_service, data_prepare_service):
        self.recommendation_service = recommendation_service
        self.data_prepare_serivice = data_prepare_service

    def GetBannerByMetaData(self, request, context):
        logger.debug(
            f"Got request platform_id={request.platform_id}, slot_name={request.slot_name}"
        )

        banner_ids = list(request.banner_id)

        try:
            platform = self.data_prepare_serivice.get_platform_by_id(
                request.platform_id
            )
            if platform == None:
                context.set_code(grpc.StatusCode.PERMISSION_DENIED)
                return banner_pb2.Banner()

            banners = self.data_prepare_serivice.get_banners(banner_ids)
            if banners == None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return banner_pb2.Banner()

            best_id = self.recommendation_service.recommend_banner(
                request.slot_name, platform.username, platform.description, banners
            )

            banner = self.data_prepare_serivice.get_proto_banner_by_id(best_id)

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
            title=banner.title,
            content=banner.content,
            description=banner.description,
            link=banner.link,
            ownerID=banner.owner_id,
            max_price=banner.max_price,
            id=best_id,
        )
