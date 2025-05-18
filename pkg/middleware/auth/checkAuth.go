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
				json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
				return
			}

			userID, role, err := authenticator.Authenticate(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
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
