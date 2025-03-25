package banner

import (
	"errors"
	"retarget/internal/banner-service/entity"
	"retarget/internal/banner-service/repo"
)

type BannerUsecase struct {
	BannerRepository *repo.BannerRepository
}

func NewBannerUsecase(bannerRepository *repo.BannerRepository) *BannerUsecase {
	return &BannerUsecase{BannerRepository: bannerRepository}
}

func (b *BannerUsecase) GetBannersByUserID(userID int) ([]*entity.Banner, error) {
	banners, err := b.BannerRepository.GetBannersByUserId(userID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerUsecase) GetBannerByID(userID, bannerID int) (*entity.Banner, error) {
	banner, err := b.BannerRepository.GetBannerByID(userID)
	if err != nil {
		return nil, err
	}
	if banner.OwnerID != userID || banner.Deleted {
		return nil, errors.New("banner not found")
	}
	return banner, err
}

func (b *BannerUsecase) DeleteBannerByID(userID, bannerID int) error {
	err := b.BannerRepository.DeleteBannerByID(userID, bannerID)
	return err
}
