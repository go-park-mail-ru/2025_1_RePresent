package csat

import (
	repo "retarget/internal/csat-service/repo/csat"
)

type CsatUsecaseInterface interface {
	ShowQuestionByPageID(page string) error
}

type CsatUsecase struct {
	csatRepository *repo.CsatRepository
}

func NewCsatUsecase(csatRepo *repo.CsatRepository) *CsatUsecase {
	return &CsatUsecase{csatRepository: csatRepo}
}

func (a *CsatUsecase) GetQuestionByPage(page string) error {
	panic("unimplimented")
}
