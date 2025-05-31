package slot

import (
	"context"
	"encoding/json"
	"net/http"

	"retarget/pkg/entity"
	response "retarget/pkg/entity"

	"retarget/internal/adv-service/usecase/slot"

	"github.com/google/uuid"
)

func (c *SlotController) DeleteSlotHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Link uuid.UUID `json:"link"`
	}
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

	if err := c.slotUsecase.DeleteSlot(context.Background(), req.Link.String(), userID); err != nil {
		response := entity.NewResponseWithBody(true, err.Error(), nil)
		if err == slot.ErrNotThisUserSlot {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response)
		return
	}

	response := entity.NewResponseWithBody(false, "Slot deleted successfully", nil)
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response)
}
