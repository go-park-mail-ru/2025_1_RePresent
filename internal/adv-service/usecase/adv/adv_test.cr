package adv

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"retarget/internal/adv-service/entity/adv"
	"retarget/internal/adv-service/mocks"
	pb "retarget/pkg/proto/banner"
	protoPayment "retarget/pkg/proto/payment"
)

// Мок для BannerServiceClient
type mockBannerServiceClient struct {
	mock.Mock
}

func (m *mockBannerServiceClient) GetRandomBanner(ctx context.Context, in *pb.BannerWithMinPrice, opts ...grpc.CallOption) (*pb.Banner, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.Banner), args.Error(1)
}

func (m *mockBannerServiceClient) GetBannerByID(ctx context.Context, in *pb.BannerRequest, opts ...grpc.CallOption) (*pb.Banner, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.Banner), args.Error(1)
}

// Мок для PaymentServiceClient
type mockPaymentServiceClient struct {
	mock.Mock
}

func (m *mockPaymentServiceClient) RegUserActivity(ctx context.Context, in *protoPayment.PaymentRequest, opts ...grpc.CallOption) (*protoPayment.PaymentResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*protoPayment.PaymentResponse), args.Error(1)
}

// Мок для AdvRepositoryInterface
type mockAdvRepository struct {
	mock.Mock
}

func (m *mockAdvRepository) FindLinksByUser(userID int) ([]adv.Link, error) {
	args := m.Called(userID)
	return args.Get(0).([]adv.Link), args.Error(1)
}

func (m *mockAdvRepository) FindUserByLink(link string) (int, error) {
	args := m.Called(link)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockAdvRepository) CreateLink(link adv.Link) error {
	args := m.Called(link)
	return args.Error(0)
}

// Добавляем недостающий метод DeleteLink
func (m *mockAdvRepository) DeleteLink(link string) error {
	args := m.Called(link)
	return args.Error(0)
}

func (m *mockAdvRepository) WriteMetric(bannerID int, slotLink, action, price string) error {
	args := m.Called(bannerID, slotLink, action, price)
	return args.Error(0)
}

func (m *mockAdvRepository) GetSlotMetric(slotLink, activity string, from, to time.Time) (map[string]int, error) {
	args := m.Called(slotLink, activity, from, to)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *mockAdvRepository) GetBannerMetric(bannerID int, activity string, from, to time.Time) (map[string]int, error) {
	args := m.Called(bannerID, activity, from, to)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *mockAdvRepository) GetSlotCTR(slotLink, activity string, from, to time.Time) (map[string]float64, error) {
	args := m.Called(slotLink, activity, from, to)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *mockAdvRepository) GetSlotAVGPrice(slotLink, activity string, from, to time.Time) (map[string]float64, error) {
	args := m.Called(slotLink, activity, from, to)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *mockAdvRepository) GetSlotRevenue(slotLink, activity string, from, to time.Time) (map[string]float64, error) {
	args := m.Called(slotLink, activity, from, to)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *mockAdvRepository) GetBannerCTR(bannerID int, activity string, from, to time.Time) (map[string]float64, error) {
	args := m.Called(bannerID, activity, from, to)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *mockAdvRepository) GetBannerExpenses(bannerID int, activity string, from, to time.Time) (map[string]float64, error) {
	args := m.Called(bannerID, activity, from, to)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func TestAdvUsecase_GetLinks(t *testing.T) {
	// Arrange
	advRepo := new(mockAdvRepository)
	bannerClient := new(mockBannerServiceClient)
	paymentClient := new(mockPaymentServiceClient)
	slotRepo := new(mocks.SlotRepositoryInterface)

	usecase := NewAdvUsecase(advRepo, bannerClient, paymentClient, slotRepo)

	userID := 1
	expectedLinks := []adv.Link{
		{
			TextLink: "link-1",
			UserID:   userID,
			Height:   100,
			Width:    200,
		},
		{
			TextLink: "link-2",
			UserID:   userID,
			Height:   300,
			Width:    400,
		},
	}

	advRepo.On("FindLinksByUser", userID).Return(expectedLinks, nil)

	// Act
	links, err := usecase.GetLinks(userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLinks, links)
	advRepo.AssertExpectations(t)
}

// ... (остальные тесты без изменений)

func TestAdvUsecase_GetLinks_Error(t *testing.T) {
	// Arrange
	advRepo := new(mockAdvRepository)
	bannerClient := new(mockBannerServiceClient)
	paymentClient := new(mockPaymentServiceClient)
	slotRepo := new(mocks.SlotRepositoryInterface)

	usecase := NewAdvUsecase(advRepo, bannerClient, paymentClient, slotRepo)

	userID := 1
	expectedErr := errors.New("database error")

	advRepo.On("FindLinksByUser", userID).Return([]adv.Link{}, expectedErr)

	// Act
	_, err := usecase.GetLinks(userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	advRepo.AssertExpectations(t)
}

// ... (остальные тесты остаются без изменений)

// Добавим тест для метода DeleteLink
func TestAdvUsecase_DeleteLink(t *testing.T) {
	// Arrange
	advRepo := new(mockAdvRepository)
	bannerClient := new(mockBannerServiceClient)
	paymentClient := new(mockPaymentServiceClient)
	slotRepo := new(mocks.SlotRepositoryInterface)

	usecase := NewAdvUsecase(advRepo, bannerClient, paymentClient, slotRepo)

	link := "link-to-delete"
	advRepo.On("DeleteLink", link).Return(nil)

	// Act
	err := usecase.DeleteLink(link)

	// Assert
	assert.NoError(t, err)
	advRepo.AssertExpectations(t)
}

func TestAdvUsecase_DeleteLink_Error(t *testing.T) {
	// Arrange
	advRepo := new(mockAdvRepository)
	bannerClient := new(mockBannerServiceClient)
	paymentClient := new(mockPaymentServiceClient)
	slotRepo := new(mocks.SlotRepositoryInterface)

	usecase := NewAdvUsecase(advRepo, bannerClient, paymentClient, slotRepo)

	link := "link-to-delete"
	expectedErr := errors.New("link not found")
	advRepo.On("DeleteLink", link).Return(expectedErr)

	// Act
	err := usecase.DeleteLink(link)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	advRepo.AssertExpectations(t)
}
