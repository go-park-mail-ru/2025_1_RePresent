package slot

import (
	"context"
	"encoding/json"
	"net/http"

	"retarget/internal/adv-service/dto"
	"retarget/internal/adv-service/usecase/slot"
	"retarget/pkg/entity"
	response "retarget/pkg/entity"
)

func (c *SlotController) EditSlotHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := entity.NewResponseWithBody(true, "Invalid request body", nil)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	updatedSlot, err := c.slotUsecase.UpdateSlot(context.Background(), req, userID)
	if err != nil {
		response := entity.NewResponseWithBody(true, err.Error(), nil)
		if err == slot.ErrNotThisUserSlot {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	responseSlot := struct {
		Link       string `json:"link"`
		SlotName   string `json:"slot_name"`
		FormatCode int    `json:"format_code"`
		MinPrice   string `json:"min_price"`
		IsActive   bool   `json:"is_active"`
	}{
		Link:       updatedSlot.Link,
		SlotName:   updatedSlot.SlotName,
		FormatCode: updatedSlot.FormatCode,
		MinPrice:   updatedSlot.MinPrice.String(),
		IsActive:   updatedSlot.IsActive,
	}

	response := entity.NewResponseWithBody(false, "Slot updated successfully", responseSlot)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
