package adv

import (
	"errors"
	"fmt"
	"regexp"
	repoAdv "retarget/internal/adv-service/repo/adv"

	"github.com/google/uuid"
)

type AdvUsecaseInterface interface {
	GetLink(user_id int) (string, error)
	CheckLink(platform_id string) error
	PutLink(user_id int) (string, bool, error)
	generateLink(user_id int) string
}

type AdvUsecase struct {
	advRepository *repoAdv.AdvRepository
}

func NewAdvUsecase(advRepo *repoAdv.AdvRepository) *AdvUsecase {
	return &AdvUsecase{advRepository: advRepo}
}

func (a *AdvUsecase) GetLink(user_id int) (string, error) {
	link, err := a.advRepository.FindLinkByUser(user_id)
	if err != nil {
		return "", err
	}
	if link == "" {
		return "", errors.New("Link by user_id not found")
	}
	return link, nil
}

func (a *AdvUsecase) CheckLink(secret_link string) error {
	if secret_link == "" {
		return errors.New("Link is empty")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(secret_link) {
		return errors.New("Invalid link")
	}
	user, err := a.advRepository.FindUserByLink(secret_link)
	if err != nil {
		return errors.New("link not found")
	}
	fmt.Println(user)
	return nil
}

func (a *AdvUsecase) PutLink(user_id int) (string, bool, error) {
	link, err := a.advRepository.FindLinkByUser(user_id)
	if err != nil {
		return "", false, err
	}
	if link != "" {
		return link, false, err
	}
	new_link := a.generateLink(user_id)
	err = a.advRepository.CreateLinkByUser(user_id, new_link)
	if err != nil {
		return "", false, err
	}
	return new_link, true, nil
}

func (a *AdvUsecase) generateLink(user_id int) string {
	return uuid.NewString()
}
