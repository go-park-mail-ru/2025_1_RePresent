package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"retarget/internal/banner-service/usecase" // Импорт usecase
	entity "retarget/pkg/entity"
	bannerpb "retarget/pkg/proto/banner" // Импорт сгенерированного gRPC-кода
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BannerServer struct {
	bannerpb.UnimplementedBannerServiceServer
	bannerUC usecase.BannerUsecase
}

// NewBannerServer — конструктор для инициализации сервера с зависимостями
func NewBannerServer(bannerUC usecase.BannerUsecase) *BannerServer {
	return &BannerServer{
		bannerUC: bannerUC,
	}
}

func (s *BannerServer) GetRandomBanner(
	ctx context.Context,
	req *bannerpb.BannerWithMinPrice,
) (*bannerpb.Banner, error) {
	dec, _ := entity.NewDec(req.MinPrice)
	banner, err := s.bannerUC.GetRandomBannerForADV(0, "", dec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get banner: %v", err)
	}
	return &bannerpb.Banner{
		Title:       banner.Title,
		Content:     banner.Content,
		Description: banner.Description,
		Link:        banner.Link,
		OwnerID:     strconv.Itoa(banner.OwnerID),
		Id:          int64(banner.ID),
		MaxPrice:    banner.MaxPrice.String(),
	}, nil
}

func (s *BannerServer) GetSuitableBanners(
	ctx context.Context,
	req *bannerpb.BannerWithMinPrice,
) (*bannerpb.ActiveBanners, error) {
	dec, _ := entity.NewDec(req.MinPrice)
	bannerIDs, err := s.bannerUC.GetSuitableBannersForADV(dec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get banner: %v", err)
	}

	return &bannerpb.ActiveBanners{
		BannerId: bannerIDs,
	}, nil
}

func (s *BannerServer) GetBannerByID(ctx context.Context, req *bannerpb.BannerRequest) (*bannerpb.Banner, error) {
	bannerID := int(req.GetId())
	banner, err := s.bannerUC.BannerRepository.GetBannerByID(bannerID, "grpc request")
	if err != nil {
		return nil, fmt.Errorf("banner is not exist")
	}

	return &bannerpb.Banner{
		Title:       banner.Title,
		Content:     banner.Content,
		Description: banner.Description,
		Link:        banner.Link,
		OwnerID:     strconv.Itoa(banner.OwnerID),
		Id:          int64(banner.ID),
		MaxPrice:    banner.MaxPrice.String(),
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
