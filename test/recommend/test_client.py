import grpc
import pkg.proto.recommend.recommend_pb2 as recommend_pb2
import pkg.proto.recommend.recommend_pb2_grpc as recommend_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50055") as channel:
        stub = recommend_pb2_grpc.RecommendServiceStub(channel)
        response = stub.GetBannerByMetaData(
            recommend_pb2.RecommendationRequest(test="Hello from client")
        )
        print("Response:", response)


if __name__ == "__main__":
    run()
