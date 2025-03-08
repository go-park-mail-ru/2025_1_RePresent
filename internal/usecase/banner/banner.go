package banner

import (
	"RE/internal/entity"
	"RE/internal/repo"
	"strconv"
)

type BannerUsecase struct {
	BannerRepository *repo.BannerRepository
}

func NewBannerUsecase(bannerRepository *repo.BannerRepository) *BannerUsecase {
	return &BannerUsecase{BannerRepository: bannerRepository}
}

// GetBannersByUserID возвращает все баннеры для пользователя по user_id
func (u *BannerUsecase) GetBannersByUserID(userIDStr string) ([]*entity.Banner, error) {
	// Преобразуем user_id из строки в int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, err
	}

	// Используем репозиторий для получения баннеров
	banners, err := u.BannerRepository.GetBannersByUserId(userID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}
