package adv

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	entity "retarget/pkg/entity"
	"strconv"

	"github.com/gorilla/mux"
)

type IFrame struct {
	ImageSrc    string
	Link        string
	Title       string
	Description string
	Banner      string
	Slot        string
}

func (c *AdvController) IframeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret_link := vars["link"]
	banner, err := c.advUsecase.GetIframe(secret_link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	data := IFrame{
		ImageSrc:    "http://re-target.ru/api/v1/banner/image/" + banner.Content,
		Link:        banner.Link,
		Title:       banner.Title,
		Description: banner.Description,
		Banner:      strconv.Itoa(int(banner.Id)),
		Slot:        secret_link,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
