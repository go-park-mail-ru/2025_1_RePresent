import grpc
import os, sys

TEST_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
sys.path.append(TEST_ROOT)


import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50055") as channel:
        stub = recommend_pb2_grpc.RecommendServiceStub(channel)
        request = recommend_pb2.RecommendationRequest(
            platform_id=100, slot_name="Хомяки", banner_id=[1, 2, 3, 4]
        )
        print("📤 Отправляем запрос...")
        try:
            response = stub.GetBannerByMetaData(request)
            print("✅ Получен ответ:")
            print(f"ID: {response.id}")
            print(f"Title: {response.title}")
            print(f"Max Price: {response.max_price}")
        except grpc.RpcError as e:
            print(f"❌ Ошибка gRPC: {e.code()} — {e.details()}")


if __name__ == "__main__":
    run()
