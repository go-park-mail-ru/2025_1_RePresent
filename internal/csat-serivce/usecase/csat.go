package csat

import (
	"net/http"
	"retarget/internal/csat-serivce/repo"
)

type CsatUsecase struct {
	CsatRepository *repo.CsatRepository
}

func (h *CsatController) ShowQuestionByPageID(w http.ResponseWriter, r *http.Request) {
}
