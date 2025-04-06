package adv

import (
	"encoding/json"
	"net/http"
	entity "retarget/pkg/entity"
)

func (c *AdvController) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got"))
}
