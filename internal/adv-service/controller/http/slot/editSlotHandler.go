package slot

import (
	"context"
	"encoding/json"
	"net/http"

	"retarget/internal/adv-service/dto"
	model "retarget/internal/adv-service/easyjsonModels"
	"retarget/internal/adv-service/usecase/slot"
	"retarget/pkg/entity"
	response "retarget/pkg/entity"
)

func (c *SlotController) EditSlotHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := entity.NewResponseWithBody(true, "Invalid request body", nil)
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response)
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
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
		//nolint:errcheck
		json.NewEncoder(w).Encode(response)
		return
	}

	responseSlot := model.EditSlotResponse{
		Link:       updatedSlot.Link,
		SlotName:   updatedSlot.SlotName,
		FormatCode: updatedSlot.FormatCode,
		MinPrice:   updatedSlot.MinPrice.String(),
		IsActive:   updatedSlot.IsActive,
	}

	response := entity.NewResponseWithBody(false, "Slot updated successfully", responseSlot)
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response)
}
