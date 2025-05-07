package adv

import (
	"encoding/json"
	"log"
	"net/http"
	entity "retarget/pkg/entity"

	"github.com/gorilla/mux"
)

func (c *AdvController) IframeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret_link := vars["link"]
	log.Println("ало работай")
	c.advUsecase.GetIframe(secret_link)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got"))
}
