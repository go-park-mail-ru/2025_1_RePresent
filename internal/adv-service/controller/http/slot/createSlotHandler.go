package slot

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"retarget/internal/adv-service/dto"
	"retarget/pkg/entity"
	response "retarget/pkg/entity"
)

func (c *SlotController) CreateSlotHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest
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

	createdSlot, err := c.slotUsecase.CreateSlot(context.Background(), req, userID)
	if err != nil {
		response := entity.NewResponseWithBody(true, err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	responseSlot := struct {
		Link       string    `json:"link"`
		SlotName   string    `json:"slot_name"`
		FormatCode int       `json:"format_code"`
		MinPrice   string    `json:"min_price"`
		IsActive   bool      `json:"is_active"`
		CreatedAt  time.Time `json:"created_at"`
	}{
		Link:       createdSlot.Link,
		SlotName:   createdSlot.SlotName,
		FormatCode: createdSlot.FormatCode,
		MinPrice:   createdSlot.MinPrice.String(),
		IsActive:   createdSlot.IsActive,
		CreatedAt:  createdSlot.CreatedAt,
	}

	response := entity.NewResponseWithBody(false, "Slot created successfully", responseSlot)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
