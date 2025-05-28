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
	"github.com/mailru/easyjson"
)

type IFrame struct {
	ImageSrc    string
	Link        string
	Title       string
	Description string
	Action      string
	Banner      int
	Slot        string
}

func (h *BannerController) GetBannerIFrameByID(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner, err := h.BannerUsecase.GetBannerForIFrame(bannerID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}
	url, err := h.LinkBuilder.BannerImageURL(banner.Content)
	if err != nil {
		log.Println("Обработка ошибки")
		// обработка ошибки
	}
	data := IFrame{
		ImageSrc:    url,
		Link:        banner.Link,
		Title:       banner.Title,
		Description: banner.Description,
		Action:      "",
		Banner:      -1,
		Slot:        "",
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *BannerController) RandomIFrame(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	vars := mux.Vars(r)
	userIDstr := vars["uniq_link"]
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		resp := response.NewResponse(false, "invalid banner ID")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	banner, err := h.BannerUsecase.GetRandomBannerForIFrame(userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		resp := response.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	url, err := h.LinkBuilder.BannerImageURL(banner.Content)
	if err != nil {
		log.Println("Обработка ошибки")
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
