package controller

import (
	"encoding/json"
	"html/template"
	"path/filepath"
	"strconv"

	// "fmt"
	"net/http"
	// Хардкод
	// "strconv"

	// pkg "retarget/internal/pkg/entity"
	response "pkg/entity"
	// "strconv"
	"github.com/gorilla/mux"
)

type IFrame struct {
	ImageSrc    string
	Link        string
	Title       string
	Description string
}

func (h *BannerController) GetBannerIFrame(w http.ResponseWriter, r *http.Request) {

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "card.html")))
	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	banner, err := h.BannerUsecase.GetBannerByID(userID, bannerID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	data := IFrame{
		ImageSrc:    "http://109.120.190.243/api/v1/banner/" + banner.Content,
		Link:        banner.Link,
		Title:       banner.Title,
		Description: banner.Description,
	}
	tmpl.Execute(w, data)
}
