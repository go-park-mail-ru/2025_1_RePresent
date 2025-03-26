package adv

import (
	"encoding/json"
	"net/http"
	entity "retarget/pkg/entity"
)

func (c *AdvController) IframeHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// secret_link := vars["link"]

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got"))
}
