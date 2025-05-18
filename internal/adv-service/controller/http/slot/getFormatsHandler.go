package slot

import (
	"context"
	"encoding/json"
	"net/http"

	"retarget/pkg/entity"
)

func (c *SlotController) GetFormatsHandler(w http.ResponseWriter, r *http.Request) {
	formats, err := c.slotUsecase.GetFormats(context.Background())
	if err != nil {
		response := entity.NewResponseWithBody(true, err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := entity.NewResponseWithBody(false, "Formats retrieved successfully", formats)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
