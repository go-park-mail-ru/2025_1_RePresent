import os
import sys
import grpc

TEST_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
sys.path.append(TEST_ROOT)

# Proto
import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50055") as channel:
        stub = recommend_pb2_grpc.RecommendServiceStub(channel)
        request = recommend_pb2.RecommendationRequest(
            platform_id=1,
            slot_name="Test Name",
            banner_id=[1, 2, 3, 4, 5],
        )
        print("Sending request:", request)
        try:
            response = stub.GetBannerByMetaData(request)
            print("Received response:")
            print(f"  ID: {response.id}")
            print(f"  Title: {response.title}")
            print(f"  Description: {response.description}")
            print(f"  Max Price: {response.max_price}")
        except grpc.RpcError as e:
            print(f"gRPC Error [{e.code()}]: {e.details()}")


if __name__ == "__main__":
    run()
