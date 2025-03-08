package banner

import (
	"RE/internal/entity"
	"RE/internal/repo"
)

type BannerUsecase struct {
	BannerRepository *repo.BannerRepository
}

func NewBannerUsecase(bannerRepository *repo.BannerRepository) *BannerUsecase {
	return &BannerUsecase{BannerRepository: bannerRepository}
}

func (u *BannerUsecase) GetBannersByUserID(userID int) ([]*entity.Banner, error) {
	banners, err := u.BannerRepository.GetBannersByUserId(userID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}
