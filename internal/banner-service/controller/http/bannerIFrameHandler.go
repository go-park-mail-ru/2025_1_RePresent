package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"path/filepath"
	"strconv"

	"net/http"
	response "retarget/pkg/entity"

	"github.com/gorilla/mux"
)

type IFrame struct {
	ImageSrc    string
	Link        string
	Title       string
	Description string
}

func (h *BannerController) GetBannerIFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner, err := h.BannerUsecase.GetBannerByID(userID, bannerID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}
	url, err := h.LinkBuilder.BannerImageURL(banner.Content)
	if err != nil {
		// обработка ошибки
	}
	data := IFrame{
		ImageSrc:    url,
		Link:        banner.Link,
		Title:       banner.Title,
		Description: banner.Description,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
