package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"path/filepath"
	"strconv"

	// "fmt"
	"net/http"
	// Хардкод
	// "strconv"

	// pkg "retarget/internal/pkg/entity"
	response "retarget/pkg/entity"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	// if !ok {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	// }
	// userID := userSession.UserID
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner, err := h.BannerUsecase.GetBannerByID(1, bannerID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := IFrame{
		ImageSrc:    "http://109.120.190.243/api/v1/banner/image/" + banner.Content,
		Link:        banner.Link,
		Title:       banner.Title,
		Description: banner.Description,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
