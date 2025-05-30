package usecase

import (
	"errors"
	"math/rand"
	model "retarget/internal/banner-service/easyjsonModels"
	"retarget/internal/banner-service/entity"
	"retarget/internal/banner-service/repo"
	"time"

	decimal "retarget/pkg/entity"
)

type BannerRepo interface {
	GetBannersByUserId(int, string) ([]model.Banner, error)
	GetBannerByID(int, string) (*model.Banner, error)
	GetMaxPriceBanner(interface{}) *model.Banner
	CreateNewBanner(model.Banner, string) error
	UpdateBanner(model.Banner, string) error
	DeleteBannerByID(int, int, string) error
}

type BannerUsecase struct {
	BannerRepository *repo.BannerRepository
	rng              *rand.Rand
}

func NewBannerUsecase(bannerRepository *repo.BannerRepository) *BannerUsecase {
	return &BannerUsecase{BannerRepository: bannerRepository, rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (b *BannerUsecase) GetBannersByUserID(userID int, requestID string) (model.BannerList, error) {
	banners, err := b.BannerRepository.GetBannersByUserId(userID, requestID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerUsecase) GetBannerByID(userID, bannerID int, requestID string) (*model.Banner, error) {
	banner, err := b.BannerRepository.GetBannerByID(bannerID, requestID)
	if err != nil {
		return nil, err
	}
	if banner.OwnerID != userID || banner.Deleted {
		return nil, errors.New("banner not found")
	}
	return banner, err
}

func (b *BannerUsecase) GetBannerForIFrame(bannerID int, requestID string) (*model.Banner, error) {
	banner, err := b.BannerRepository.GetBannerByID(bannerID, requestID)
	if err != nil {
		return nil, err
	}
	return banner, err
}

func (b *BannerUsecase) GetRandomBannerForIFrame(userID int, requestID string) (*model.Banner, error) {
	banners, err := b.BannerRepository.GetBannersByUserId(userID, requestID)
	if err != nil {
		return nil, err
	}
	if len(banners) > 0 {
		return &banners[b.rng.Intn(len(banners))], nil
	}
	return &entity.DefaultBanner, nil
}

func (b *BannerUsecase) GetRandomBannerForADV(userID int, requestID string, floor *decimal.Decimal) (*model.Banner, error) {
	banner := b.BannerRepository.GetMaxPriceBanner(floor)
	if banner == nil {
		return &entity.DefaultBanner, nil
	}
	return banner, nil
}

func (b *BannerUsecase) GetSuitableBannersForADV(floor *decimal.Decimal) ([]int64, error) {
	bannerIDs, err := b.BannerRepository.GetSuitableBanners(floor)
	if err != nil {
		return []int64{-1}, nil
	}
	if len(bannerIDs) == 1 && bannerIDs[0] == -1 {
		return []int64{-1}, nil
	}
	return bannerIDs, nil
}

func (b *BannerUsecase) CreateBanner(userID int, banner model.Banner, requestID string) error {
	err := b.BannerRepository.CreateNewBanner(banner, requestID)
	return err
}

func (b *BannerUsecase) UpdateBanner(userID int, banner model.Banner, requestID string) error {
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

func (b *BannerUsecase) GenerateBannerDescription(userID, bannerID int, requestID string) (string, error) {
	banner, err := b.BannerRepository.GetBannerByID(bannerID, requestID)
	if err != nil {
		return "", err
	}

	if banner.OwnerID != userID || banner.Deleted {
		return "", errors.New("banner not found or access denied")
	}

	description, err := b.BannerRepository.GenerateBannerDescription(bannerID, requestID)
	if err != nil {
		return "", err
	}

	banner.Description = description
	err = b.BannerRepository.UpdateBanner(*banner, requestID)
	if err != nil {
		return "", err
	}

	return description, nil
}

func (b *BannerUsecase) GenerateBannerImage(userID, bannerID int, requestID string) (string, error) {
	banner, err := b.BannerRepository.GetBannerByID(bannerID, requestID)
	if err != nil {
		return "", err
	}
	if banner.OwnerID != userID || banner.Deleted {
		return "", errors.New("banner not found or access denied")
	}
	return b.BannerRepository.GenerateBannerImage(bannerID, requestID)
}
