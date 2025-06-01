package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"retarget/pkg/entity"
)

func AuthMiddleware(authenticator AuthenticatorInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)

				if errEncode := json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error())); errEncode != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}

				return
			}

			userID, role, err := authenticator.Authenticate(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)

				encodeErr := json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
				if encodeErr != nil {
					http.Error(w, "Failed to write response", http.StatusInternalServerError)
					return
				}

				return
			}

			userContext := entity.UserContext{
				UserID: userID,
				Role:   role,
			}

			ctx := context.WithValue(r.Context(), entity.UserContextKey, userContext)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
