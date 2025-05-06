package server

import (
	"context"
	"log"
	"net"

	"retarget/internal/banner-service/usecase" // Импорт usecase
	bannerpb "retarget/pkg/proto"              // Импорт сгенерированного gRPC-кода

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BannerServer struct {
	bannerpb.UnimplementedBannerServiceServer
	bannerUC usecase.BannerUsecase // Добавляем зависимость
}

// NewBannerServer — конструктор для инициализации сервера с зависимостями
func NewBannerServer(bannerUC usecase.BannerUsecase) *BannerServer {
	return &BannerServer{
		bannerUC: bannerUC,
	}
}

func (s *BannerServer) GetRandomBanner(
	ctx context.Context,
	req *bannerpb.Empty,
) (*bannerpb.Banner, error) {
	banner, err := s.bannerUC.GetRandomBannerForADV(0, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get banner: %v", err)
	}

	return &bannerpb.Banner{
		Title:       banner.Title,
		Content:     banner.Content,
		Description: banner.Description,
		Link:        banner.Link,
	}, nil
}

func RunGRPCServer(bannerUC usecase.BannerUsecase) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем сервер и передаем ему usecase
	s := grpc.NewServer()
	bannerpb.RegisterBannerServiceServer(s, NewBannerServer(bannerUC))

	log.Printf("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
