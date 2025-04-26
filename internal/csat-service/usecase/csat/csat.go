package csat

import (
	"fmt"
	"math/rand"
	entity "retarget/internal/csat-service/entity/csat"
	repo "retarget/internal/csat-service/repo/csat"
)

type CsatUsecaseInterface interface {
	ShowQuestionByPageID(page string) (string, error)
	CreateReview(page string) error
}

type CsatUsecase struct {
	csatRepository *repo.CsatRepository
}

func NewCsatUsecase(csatRepo *repo.CsatRepository) *CsatUsecase {
	return &CsatUsecase{csatRepository: csatRepo}
}

func (a *CsatUsecase) GetQuestionByPage(page string) (string, error) {
	res, err := a.csatRepository.GetQuestionsByPage(page)
	if err != nil {
		return " ", fmt.Errorf("Error fetching question")
	}
	if len(res) > 0 {
		return res[rand.Intn(len(res))], nil
	}
	return " ", fmt.Errorf("Question not found")
}

func (a *CsatUsecase) CreateReview(review entity.Review) error {
	err := a.csatRepository.AddReview(review)
	if err != nil {
		return err
	}
	return nil
}
