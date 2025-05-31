package slot

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/inf.v0"

	"retarget/internal/adv-service/dto"
	"retarget/internal/adv-service/entity/slot"
	"retarget/internal/adv-service/mocks"
	repoSlot "retarget/internal/adv-service/repo/slot"
)

func TestSlotUsecase_CreateSlot(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	req := dto.CreateRequest{
		SlotName:   "Test Slot",
		FormatCode: 1,
		MinPrice:   *inf.NewDec(100, 0),
		IsActive:   true,
	}

	expectedSlot := slot.Slot{
		Link:       "mocked-uuid",
		SlotName:   req.SlotName,
		FormatCode: req.FormatCode,
		MinPrice:   req.MinPrice,
		IsActive:   req.IsActive,
		CreatedAt:  time.Now().UTC(),
	}

	// Need to use mock.AnythingOfType for Slot because CreatedAt will be different
	repoMock.On("CreateSlot", mock.Anything, userID, mock.AnythingOfType("slot.Slot")).Return(expectedSlot, nil)

	// Act
	result, err := uc.CreateSlot(context.Background(), req, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedSlot, result)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_CheckLink(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	link := "valid-link"
	repoMock.On("GetUserByLink", mock.Anything, link).Return(1, time.Now(), nil)

	// Act
	err := uc.CheckLink(link)

	// Assert
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_CheckLink_Error(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	link := "invalid-link"
	expectedErr := repoSlot.ErrSlotNotFound
	repoMock.On("GetUserByLink", mock.Anything, link).Return(0, time.Time{}, expectedErr)

	// Act
	err := uc.CheckLink(link)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_UpdateSlot(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	link := uuid.New()
	createdAt := time.Now().UTC()

	req := dto.UpdateRequest{
		Link:       link,
		SlotName:   "Updated Slot",
		FormatCode: 2,
		MinPrice:   *inf.NewDec(200, 0),
		IsActive:   false,
	}

	expectedSlot := slot.Slot{
		Link:       link.String(),
		SlotName:   req.SlotName,
		FormatCode: req.FormatCode,
		MinPrice:   req.MinPrice,
		IsActive:   req.IsActive,
		CreatedAt:  createdAt,
	}

	repoMock.On("GetUserByLink", mock.Anything, link.String()).Return(userID, createdAt, nil)
	repoMock.On("UpdateSlot", mock.Anything, userID, mock.MatchedBy(func(s slot.Slot) bool {
		return s.Link == link.String() && s.CreatedAt == createdAt
	})).Return(nil)

	// Act
	result, err := uc.UpdateSlot(context.Background(), req, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedSlot, result)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_UpdateSlot_NotFound(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	link := uuid.New()

	req := dto.UpdateRequest{
		Link:       link,
		SlotName:   "Updated Slot",
		FormatCode: 2,
		MinPrice:   *inf.NewDec(200, 0),
		IsActive:   false,
	}

	repoMock.On("GetUserByLink", mock.Anything, link.String()).Return(0, time.Time{}, repoSlot.ErrSlotNotFound)

	// Act
	_, err := uc.UpdateSlot(context.Background(), req, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrSlotNotFound, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_UpdateSlot_NotThisUserSlot(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	otherUserID := 2
	link := uuid.New()
	createdAt := time.Now().UTC()

	req := dto.UpdateRequest{
		Link:       link,
		SlotName:   "Updated Slot",
		FormatCode: 2,
		MinPrice:   *inf.NewDec(200, 0),
		IsActive:   false,
	}

	repoMock.On("GetUserByLink", mock.Anything, link.String()).Return(otherUserID, createdAt, nil)

	// Act
	_, err := uc.UpdateSlot(context.Background(), req, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNotThisUserSlot, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_DeleteSlot(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	link := "valid-link"
	createdAt := time.Now().UTC()

	repoMock.On("GetUserByLink", mock.Anything, link).Return(userID, createdAt, nil)
	repoMock.On("DeleteSlot", mock.Anything, userID, link, createdAt).Return(nil)

	// Act
	err := uc.DeleteSlot(context.Background(), link, userID)

	// Assert
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_DeleteSlot_NotFound(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	link := "invalid-link"
	expectedErr := repoSlot.ErrSlotNotFound

	repoMock.On("GetUserByLink", mock.Anything, link).Return(0, time.Time{}, expectedErr)

	// Act
	err := uc.DeleteSlot(context.Background(), link, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_DeleteSlot_NotThisUserSlot(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	otherUserID := 2
	link := "valid-link"
	createdAt := time.Now().UTC()

	repoMock.On("GetUserByLink", mock.Anything, link).Return(otherUserID, createdAt, nil)

	// Act
	err := uc.DeleteSlot(context.Background(), link, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrNotThisUserSlot, err)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_GetUserSlots(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	userID := 1
	expectedSlots := []slot.Slot{
		{
			Link:       "link-1",
			SlotName:   "Slot 1",
			FormatCode: 1,
			MinPrice:   *inf.NewDec(100, 0),
			IsActive:   true,
			CreatedAt:  time.Now().UTC(),
		},
		{
			Link:       "link-2",
			SlotName:   "Slot 2",
			FormatCode: 2,
			MinPrice:   *inf.NewDec(200, 0),
			IsActive:   false,
			CreatedAt:  time.Now().UTC(),
		},
	}

	repoMock.On("GetSlotsByUser", mock.Anything, userID).Return(expectedSlots, nil)

	// Act
	result, err := uc.GetUserSlots(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedSlots, result)
	repoMock.AssertExpectations(t)
}

func TestSlotUsecase_GetFormats(t *testing.T) {
	// Arrange
	repoMock := new(mocks.SlotRepositoryInterface)
	uc := NewSlotUsecase(repoMock)

	expectedFormats := []slot.Format{
		{
			Code:        1,
			Height:      100,
			Width:       200,
			Description: "Banner 100x200",
		},
		{
			Code:        2,
			Height:      300,
			Width:       400,
			Description: "Banner 300x400",
		},
	}

	repoMock.On("GetCurrentFormats", mock.Anything).Return(expectedFormats, nil)

	// Act
	result, err := uc.GetFormats(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedFormats, result)
	repoMock.AssertExpectations(t)
}
