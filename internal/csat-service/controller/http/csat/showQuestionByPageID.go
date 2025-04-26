package csat

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	// "retarget/csat-service/entity"
	response "retarget/pkg/entity"
	// "strconv"

	"github.com/gorilla/mux"
)

type IFrame struct {
	Question string
}

func (c *CsatController) ShowQuestionByPageID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler")
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
	json.NewEncoder(w).Encode(response.NewResponseWithBody(false, "Success send", question))
}

func (c *CsatController) ShowQuestionIFrameByPageID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler")
	_, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "iframe.html")))
	vars := mux.Vars(r)
	pageID := vars["page_id"]
	question, err := c.csatUsecase.GetQuestionByPage(pageID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}
	data := IFrame{
		Question: question,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
