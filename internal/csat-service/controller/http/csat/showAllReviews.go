package csat

import (
	"encoding/json"
	"net/http"

	// "retarget/csat-service/entity"
	response "retarget/pkg/entity"
)

func (c *CsatController) ShowReviewsByUserID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		return
	}
	// userID := userSession.UserID
	reviews, err := c.csatUsecase.GetAllReviews()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewResponseWithBody(false, "Success got", reviews))
}
