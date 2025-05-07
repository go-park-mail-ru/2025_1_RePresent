package adv

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	entity "retarget/pkg/entity"
	protoPayment "retarget/pkg/proto/payment"
	"strconv"

	"github.com/gorilla/mux"
)

type IFrame struct {
	ImageSrc    string
	Link        string
	Title       string
	Description string
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
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	ctx := context.Background()
	owner, _ := strconv.Atoi(banner.OwnerID)
	userID, _, _ := c.advUsecase.SlotsRepository.GetUserByLink(ctx, secret_link)
	request := protoPayment.PaymentRequest{FromUserId: int32(owner), ToUserId: int32(userID), Amount: 1}
	c.advUsecase.PaymentClient.RegUserActivity(ctx, &request)
}
