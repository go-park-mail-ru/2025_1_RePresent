package slot

import (
	"context"
	"errors"
	"fmt"
	"retarget/internal/adv-service/dto"
	"retarget/internal/adv-service/entity/slot"
	"time"

	repoSlot "retarget/internal/adv-service/repo/slot"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	ErrNotThisUserSlot = errors.New("This is not your slot!")
	ErrValidation      = errors.New("validation error")
	ErrSlotNotFound    = errors.New("slot not found")
	ErrAccessDenied    = errors.New("access denied")
)

type SlotUsecaseInterface interface {
	CreateSlot(ctx context.Context, req dto.CreateRequest, userID int) (slot.Slot, error)
	UpdateSlot(ctx context.Context, req dto.UpdateRequest, userID int) (slot.Slot, error)
	DeleteSlot(ctx context.Context, link string, userID int) error
	CheckLink(link string) error
	GetUserSlots(ctx context.Context, userID int) ([]slot.Slot, error)
	GetFormats(ctx context.Context) ([]slot.Format, error)
}

type SlotUsecase struct {
	repo     repoSlot.SlotRepositoryInterface
	validate *validator.Validate
}

func NewSlotUsecase(repo repoSlot.SlotRepositoryInterface) *SlotUsecase {
	return &SlotUsecase{
		repo:     repo,
		validate: validator.New(),
	}
}

func (uc *SlotUsecase) CreateSlot(ctx context.Context, req dto.CreateRequest, userID int) (slot.Slot, error) {
	s := slot.Slot{
		Link:       uuid.New().String(),
		SlotName:   req.SlotName,
		FormatCode: req.FormatCode,
		MinPrice:   req.MinPrice,
		IsActive:   req.IsActive,
		CreatedAt:  time.Now().UTC(),
	}
	return uc.repo.CreateSlot(ctx, userID, s)
}

func (uc *SlotUsecase) CheckLink(link string) error {
	ctx := context.Background()
	_, _, err := uc.repo.GetUserByLink(ctx, link)
	return err
}

func (uc *SlotUsecase) UpdateSlot(ctx context.Context, req dto.UpdateRequest, userID int) (slot.Slot, error) {
	user_id, created_at, err := uc.repo.GetUserByLink(ctx, req.Link.String())
	if err != nil {
		return slot.Slot{}, ErrSlotNotFound
	}
	if user_id != userID {
		return slot.Slot{}, ErrNotThisUserSlot
	}

	s := slot.Slot{
		Link:       req.Link.String(),
		SlotName:   req.SlotName,
		FormatCode: req.FormatCode,
		MinPrice:   req.MinPrice,
		IsActive:   req.IsActive,
		CreatedAt:  created_at,
	}

	if err := uc.repo.UpdateSlot(ctx, userID, s); err != nil {
		return slot.Slot{}, err
	}

	return s, nil
}

func (uc *SlotUsecase) DeleteSlot(ctx context.Context, link string, userID int) error {
	user_id, created_at, err := uc.repo.GetUserByLink(ctx, link)
	if err != nil {
		return err
	}
	if user_id != userID {
		return ErrNotThisUserSlot
	}
	fmt.Println(user_id)
	fmt.Println(userID)
	return uc.repo.DeleteSlot(ctx, userID, link, created_at)
}

func (uc *SlotUsecase) GetUserSlots(ctx context.Context, userID int) ([]slot.Slot, error) {
	return uc.repo.GetSlotsByUser(ctx, userID)
}

func (uc *SlotUsecase) GetFormats(ctx context.Context) ([]slot.Format, error) {
	return uc.repo.GetCurrentFormats(ctx)
}
