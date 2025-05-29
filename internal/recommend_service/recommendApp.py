from concurrent.futures import ThreadPoolExecutor
import grpc
import logging

logging.basicConfig(level=logging.INFO)

import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc
import pkg.proto.banner.banner_pb2 as banner_pb2


class RecommendService(recommend_pb2_grpc.RecommendServiceServicer):
    def GetBannerByMetaData(self, request, context):
        logging.info(f"📩 Получено сообщение от Go: {request.test}")

        return banner_pb2.Banner(
            title="Рекламный баннер",
            content="Купи слона!",
            description="Это тестовый баннер для рекомендаций",
            link="https://example.com ",
            ownerID="owner_12345",
            max_price="100.50",
            id=1,
        )


def serve():
    logging.info("Creating gRPC server")
    server = grpc.server(ThreadPoolExecutor(max_workers=10))
    recommend_pb2_grpc.add_RecommendServiceServicer_to_server(
        RecommendService(), server
    )

    server.add_insecure_port("[::]:50055")
    logging.info("gRPC Recommend Server Started on ReTargetApiRecommend:50055")
    server.start()

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logging.info("Stoping Recommend Server...")
        server.stop(grace=5)
        logging.info("Recommend Server stopped")
