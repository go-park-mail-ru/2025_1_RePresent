package csat

import (
	"encoding/json"
	"net/http"
	"time"

	entity "retarget/internal/csat-service/entity/csat"
	response "retarget/pkg/entity"
	"retarget/pkg/utils/validator"
)

type ReviewRequest struct {
	Page     string `json:"page_id" validate:"required"`
	Question string `json:"question" validate:"required"`
	Rating   int    `json:"rating" validate:"required,gte=0,lte=10"`
	Comment  string `json:"comment" validate:"lte=200"`
}

func (c *CsatController) SendReview(w http.ResponseWriter, r *http.Request) {
	// requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	var req ReviewRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewResponse(true, validate_errors))
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	review := entity.Review{
		UserID:    userID,
		Question:  req.Question,
		Rating:    req.Rating,
		Comment:   req.Comment,
		Page:      req.Page,
		CreatedAt: time.Now(),
	}

	if err := c.csatUsecase.CreateReview(review); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.NewResponse(false, "Thanks for your review"))
}
