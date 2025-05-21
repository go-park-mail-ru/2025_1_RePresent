package adv

import (
	"encoding/json"
	"net/http"
	entity "retarget/pkg/entity"

	"strconv"
)

func (c *AdvController) MetricsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	action := query.Get("action")
	bannerIDstr := query.Get("banner")
	slot := query.Get("slot")

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Format"))
		return
	}

	if err = c.advUsecase.WriteMetric(bannerID, slot, action); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got"))
}
