package adv

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/pkg/entity"
	response "retarget/pkg/entity"
)

type GenerateLinkRequest struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

func (c *AdvController) GenerateLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	// requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	var req GenerateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponseWithBody(true, "Invalid request body", nil))
		return
	}

	link, isCreated, err := c.advUsecase.PutLink(userID, req.Height, req.Width)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponseWithBody(true, err.Error(), nil))
		return
	}

	fullLink := fmt.Sprintf("http://re-target.ru/api/v1/adv/iframe/%s", link.TextLink)

	responseBody := struct {
		Link   string `json:"link"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	}{
		Link:   fullLink,
		Height: link.Height,
		Width:  link.Width,
	}

	if isCreated {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "Link created", responseBody))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "Link already exists", responseBody))
	}
}
