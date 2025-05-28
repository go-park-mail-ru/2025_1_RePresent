package slot

import (
	"context"
	"encoding/json"
	"net/http"

	"retarget/pkg/entity"
	response "retarget/pkg/entity"

	model "retarget/internal/adv-service/easyjsonModels"

	"github.com/mailru/easyjson"
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

	responseSlots := make([]model.GetSlotResponse, len(slots))

	for i, s := range slots {
		responseSlots[i] = model.GetSlotResponse{
			Link:       s.Link,
			SlotName:   s.SlotName,
			FormatCode: s.FormatCode,
			MinPrice:   s.MinPrice,
			IsActive:   s.IsActive,
			CreatedAt:  s.CreatedAt,
		}
	}
	// response := entity.NewResponseWithBody(false, "User slots retrieved successfully", responseSlots)

	s := "User slots retrieved successfully"
	serv := entity.ServiceResponse{
		Success: &s,
	}
	response := model.ResponseWithSlots{
		Service: serv,
		Body:    responseSlots,
	}
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(response)
	easyjson.MarshalToWriter(&response, w)
}
