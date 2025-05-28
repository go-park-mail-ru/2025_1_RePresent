package adv

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	model "retarget/internal/adv-service/easyjsonModels"
	entity "retarget/pkg/entity"

	"github.com/gorilla/mux"
)

func (c *AdvController) IframeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret_link := vars["link"]
	banner, err := c.advUsecase.GetIframe(secret_link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error())); encodeErr != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	data := model.IFrame{
		ImageSrc:    "https://test.re-target.ru/api/v1/banner/image/" + banner.Content,
		Link:        banner.Link,
		Title:       banner.Title,
		Description: banner.Description,
		Banner:      banner.Id,
		Slot:        secret_link,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
