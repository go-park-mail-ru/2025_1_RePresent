package banner

import (
	"RE/internal/usecase/banner"
	"net/http"
)

// SetupBannerRoutes настраивает маршруты для работы с баннерами
func SetupBannerRoutes(bannerUsecase *banner.BannerUsecase) http.Handler {
	mux := http.NewServeMux()

	// Создаем обработчик
	bannerHandler := NewBannerHandler(bannerUsecase)

	// Настроим маршрут для получения баннеров пользователя по user_id
	mux.HandleFunc("/user/{user_id}/banners", bannerHandler.GetBannersByUserID) //.Methods("GET")

	return mux
}
