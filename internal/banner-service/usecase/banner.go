package usecase

import (
	"errors"
	"math/rand"
	"retarget/internal/banner-service/entity"
	"retarget/internal/banner-service/repo"
	"time"
)

type BannerUsecase struct {
	BannerRepository *repo.BannerRepository
}

func NewBannerUsecase(bannerRepository *repo.BannerRepository) *BannerUsecase {
	return &BannerUsecase{BannerRepository: bannerRepository}
}

func (b *BannerUsecase) GetBannersByUserID(userID int, requestID string) ([]entity.Banner, error) {
	banners, err := b.BannerRepository.GetBannersByUserId(userID, requestID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerUsecase) GetBannerByID(userID, bannerID int, requestID string) (*entity.Banner, error) {
	banner, err := b.BannerRepository.GetBannerByID(bannerID, requestID)
	if err != nil {
		return nil, err
	}
	if banner.OwnerID != userID || banner.Deleted {
		return nil, errors.New("banner not found")
	}
	return banner, err
}

func (b *BannerUsecase) GetBannerForIFrame(bannerID int, requestID string) (*entity.Banner, error) {
	banner, err := b.BannerRepository.GetBannerByID(bannerID, requestID)
	if err != nil {
		return nil, err
	}
	return banner, err
}

func (b *BannerUsecase) GetRandomBannerForIFrame(userID int, requestID string) (*entity.Banner, error) {
	rand.Seed(time.Now().UnixNano())
	banners, err := b.BannerRepository.GetBannersByUserId(userID, requestID)
	if err != nil {
		return nil, err
	}
	if len(banners) > 0 {
		return &banners[rand.Intn(len(banners))], err
	}
	return &entity.DefaultBanner, nil
}

func (b *BannerUsecase) CreateBanner(userID int, banner entity.Banner, requestID string) error {
	err := b.BannerRepository.CreateNewBanner(banner, requestID)
	return err
}

func (b *BannerUsecase) UpdateBanner(userID int, banner entity.Banner, requestID string) error {
	oldBanner, err := b.BannerRepository.GetBannerByID(banner.ID, requestID)
	if err != nil {
		return err
	}
	if oldBanner.OwnerID != userID {
		return errors.New("banner not Found")
	}
	err = b.BannerRepository.UpdateBanner(banner, requestID)
	return err
}

func (b *BannerUsecase) DeleteBannerByID(userID, bannerID int, requestID string) error {
	err := b.BannerRepository.DeleteBannerByID(userID, bannerID, requestID)
	return err
}
