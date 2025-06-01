package slot

import (
	"context"
	"encoding/json"
	"net/http"

	"retarget/internal/adv-service/dto"
	model "retarget/internal/adv-service/easyjsonModels"
	"retarget/pkg/entity"
	response "retarget/pkg/entity"

	"github.com/mailru/easyjson"
)

func (c *SlotController) CreateSlotHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest
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

	createdSlot, err := c.slotUsecase.CreateSlot(context.Background(), req, userID)
	if err != nil {
		//nolint:errcheck
		response := entity.NewResponseWithBody(true, err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response)
		return
	}

	responseSlot := model.CreateSlotResponse{
		Link:       createdSlot.Link,
		SlotName:   createdSlot.SlotName,
		FormatCode: createdSlot.FormatCode,
		MinPrice:   createdSlot.MinPrice.String(),
		IsActive:   createdSlot.IsActive,
		CreatedAt:  createdSlot.CreatedAt,
	}
	s := "Slot created successfully"
	serv := entity.ServiceResponse{
		Success: &s,
	}
	response := model.ResponseWithSlot{
		Service: serv,
		Body:    responseSlot,
	}
	w.WriteHeader(http.StatusCreated)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(response)
	easyjson.MarshalToWriter(&response, w)
}
