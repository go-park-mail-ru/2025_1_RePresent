package csat

import (
	"encoding/json"
	"net/http"

	// "retarget/csat-service/entity"
	response "retarget/pkg/entity"
	// "strconv"

	"github.com/gorilla/mux"
)

func (c *CsatController) ShowQuestionByPageID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		return
	}
	vars := mux.Vars(r)
	pageID := vars["page_id"]
	question, err := c.csatUsecase.GetQuestionByPage(pageID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewResponseWithBody(false, "Error of authenticator", question))
}
