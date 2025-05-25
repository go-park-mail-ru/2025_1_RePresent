package adv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"retarget/internal/adv-service/entity/adv"
	repoAdv "retarget/internal/adv-service/repo/adv"
	repoSlots "retarget/internal/adv-service/repo/slot"
	entity "retarget/pkg/entity"
	pb "retarget/pkg/proto/banner"
	protoPayment "retarget/pkg/proto/payment"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AdvUsecaseInterface interface {
	GetLinks(userID int) ([]adv.Link, error)
	CheckLink(link string) error
	PutLink(userID int, height, width int) (adv.Link, bool, error)
	generateLink(userID int, height, width int) adv.Link
	WriteMetric(bannerID int, slotLink string, metric string) error
}

type AdvUsecase struct {
	SlotsRepository repoSlots.SlotRepositoryInterface
	advRepository   repoAdv.AdvRepositoryInterface
	bannerClient    pb.BannerServiceClient
	PaymentClient   protoPayment.PaymentServiceClient
}

func NewAdvUsecase(advRepo repoAdv.AdvRepositoryInterface, bannerClient pb.BannerServiceClient, paymentClient protoPayment.PaymentServiceClient, slotsRepository repoSlots.SlotRepositoryInterface) *AdvUsecase {
	return &AdvUsecase{
		advRepository:   advRepo,
		bannerClient:    bannerClient,
		PaymentClient:   paymentClient,
		SlotsRepository: slotsRepository,
	}
}

func (a *AdvUsecase) GetLinks(userID int) ([]adv.Link, error) {
	links, err := a.advRepository.FindLinksByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get links: %w", err)
	}

	if len(links) == 0 {
		return nil, errors.New("links not found for this user")
	}

	return links, nil
}

func (a *AdvUsecase) GetIframe(key string) (*pb.Banner, error) {
	slot, err := a.SlotsRepository.GetSlotInfoByLink(context.Background(), key)
	ownerID := strconv.Itoa(entity.DefaultBanner.OwnerID)
	defaultBanner := &pb.Banner{
		Title:       entity.DefaultBanner.Title,
		Content:     entity.DefaultBanner.Content,
		Description: entity.DefaultBanner.Description,
		Link:        entity.DefaultBanner.Link,
		OwnerID:     ownerID,
		MaxPrice:    entity.DefaultBanner.MaxPrice.String(),
		Id:          int64(entity.DefaultBanner.ID),
	}
	if err != nil {
		return defaultBanner, nil
	}
	req := &pb.BannerWithMinPrice{MinPrice: slot.MinPrice.String()}
	ctx := context.Background() // однажды мы прокинем нормально контекст, но не сегодня
	banner, err := a.bannerClient.GetRandomBanner(ctx, req)
	if err != nil {
		return defaultBanner, nil
	}
	return banner, nil
}

func (a *AdvUsecase) CheckLink(link string) error {
	if link == "" {
		return errors.New("link is empty")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(link) {
		return errors.New("invalid link format")
	}

	userID, err := a.advRepository.FindUserByLink(link)
	if err != nil {
		return fmt.Errorf("link verification failed: %w", err)
	}

	fmt.Printf("Link belongs to user ID: %d\n", userID)
	return nil
}

func (a *AdvUsecase) PutLink(userID int, height, width int) (adv.Link, bool, error) {
	existingLinks, err := a.advRepository.FindLinksByUser(userID)
	if err != nil {
		return adv.Link{}, false, fmt.Errorf("failed to check existing links: %w", err)
	}

	newLink := a.generateLink(userID, height, width)

	for _, link := range existingLinks {
		if link.TextLink == newLink.TextLink {
			return link, false, nil
		}
	}

	err = a.advRepository.CreateLink(newLink)
	if err != nil {
		return adv.Link{}, false, fmt.Errorf("failed to create link: %w", err)
	}

	return newLink, true, nil
}

func (a *AdvUsecase) generateLink(userID int, height, width int) adv.Link {
	return adv.Link{
		TextLink: uuid.NewString(),
		UserID:   userID,
		Height:   height,
		Width:    width,
	}
}

func (a *AdvUsecase) WriteMetric(bannerID int, slotLink string, action string) error {

	ownerSlotID, _, err := a.SlotsRepository.GetUserByLink(context.Background(), slotLink)
	if err != nil {
		return err
	}
	bannerReq := &pb.BannerRequest{Id: int64(bannerID)}
	ctx := context.Background() // однажды мы прокинем нормально контекст, но не сегодня
	banner, err := a.bannerClient.GetBannerByID(ctx, bannerReq)
	if err != nil {
		return fmt.Errorf("get banner error")
	}
	bannerOwnerID, err := strconv.Atoi(banner.OwnerID)
	if err != nil {
		return fmt.Errorf("get banner error")
	}
	req := &protoPayment.PaymentRequest{
		FromUserId: int32(bannerOwnerID),
		ToUserId:   int32(ownerSlotID),
		Amount:     string(banner.MaxPrice),
	}
	if err := a.advRepository.WriteMetric(bannerID, slotLink, action); err != nil {
		log.Printf("Failed to write metric: %v", err)
	}
	_, err = a.PaymentClient.RegUserActivity(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register user activity: %w", err)
	}
	return nil
}

func (a *AdvUsecase) GetSlotMetric(slotLink, activity string, userID int, from, to time.Time) (map[string]int, error) {

	ownerSlotID, _, err := a.SlotsRepository.GetUserByLink(context.Background(), slotLink)
	if err != nil || userID != ownerSlotID {
		return nil, fmt.Errorf("slot not found")
	}

	total, err := a.advRepository.GetSlotMetric(slotLink, activity, from, to)
	if err != nil {
		return nil, fmt.Errorf("slot not found")
	}

	return total, nil
}

func (a *AdvUsecase) GetBannerMetric(bannerID int, activity string, userID int, from, to time.Time) (map[string]int, error) {

	bannerReq := &pb.BannerRequest{Id: int64(bannerID)}
	ctx := context.Background() // однажды мы прокинем нормально контекст, но не сегодня
	banner, err := a.bannerClient.GetBannerByID(ctx, bannerReq)
	if err != nil {
		return nil, fmt.Errorf("banner not found")
	}
	ownerID, err := strconv.Atoi(banner.OwnerID)
	if err != nil || ownerID != userID {
		return nil, fmt.Errorf("banner not found")
	}

	total, err := a.advRepository.GetBannerMetric(bannerID, activity, from, to)
	if err != nil {
		return nil, fmt.Errorf("banner not found")
	}

	return total, nil
}
