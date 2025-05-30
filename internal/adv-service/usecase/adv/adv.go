package adv

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"retarget/internal/adv-service/entity/adv"
	repoAdv "retarget/internal/adv-service/repo/adv"
	repoSlots "retarget/internal/adv-service/repo/slot"
	pb "retarget/pkg/proto"
	protoPayment "retarget/pkg/proto/payment"

	"github.com/google/uuid"
)

type AdvUsecaseInterface interface {
	GetLinks(userID int) ([]adv.Link, error)
	CheckLink(link string) error
	PutLink(userID int, height, width int) (adv.Link, bool, error)
	generateLink(userID int, height, width int) adv.Link
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
	emptyReq := &pb.Empty{}
	ctx := context.Background() // однажды мы прокинем нормально контекст, но не сегодня
	banner, err := a.bannerClient.GetRandomBanner(ctx, emptyReq)
	if err != nil {
		return nil, err
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
