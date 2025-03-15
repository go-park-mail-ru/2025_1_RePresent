package middleware

import (
	"net/http"
	"retarget/internal/auth-service/usecase/auth"
)

func AuthMiddleware(authUsecase *auth.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, err = authUsecase.GetUserBySessionID(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
