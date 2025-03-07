package auth

import (
	"net/http"
	"time"
)

func (c *AuthController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}

	expires := time.Unix(0, 0)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  expires,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
		MaxAge:   -1,
	})

	err = c.usecase.Logout(cookie.Value)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logout successful"))
}
