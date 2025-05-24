package slot

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"retarget/pkg/entity"
	response "retarget/pkg/entity"

	"gopkg.in/inf.v0"
)

func (c *SlotController) GetUserSlotsHandler(w http.ResponseWriter, r *http.Request) {
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	slots, err := c.slotUsecase.GetUserSlots(context.Background(), userID)
	if err != nil {
		response := entity.NewResponseWithBody(true, err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response)
		return
	}

	responseSlots := make([]struct {
		Link       string    `json:"link"`
		SlotName   string    `json:"slot_name"`
		FormatCode int       `json:"format_code"`
		MinPrice   inf.Dec   `json:"min_price"`
		IsActive   bool      `json:"is_active"`
		CreatedAt  time.Time `json:"created_at"`
	}, len(slots))

	for i, s := range slots {
		responseSlots[i] = struct {
			Link       string    `json:"link"`
			SlotName   string    `json:"slot_name"`
			FormatCode int       `json:"format_code"`
			MinPrice   inf.Dec   `json:"min_price"`
			IsActive   bool      `json:"is_active"`
			CreatedAt  time.Time `json:"created_at"`
		}{
			Link:       s.Link,
			SlotName:   s.SlotName,
			FormatCode: s.FormatCode,
			MinPrice:   s.MinPrice,
			IsActive:   s.IsActive,
			CreatedAt:  s.CreatedAt,
		}
	}
	response := entity.NewResponseWithBody(false, "User slots retrieved successfully", responseSlots)
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response)
}
