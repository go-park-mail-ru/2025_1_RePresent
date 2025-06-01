package server

import (
	"context"
	"log"
	"net"
	"strconv"

	entity "retarget/internal/pay-service/entity"
	usecase "retarget/internal/pay-service/usecase" // Импорт usecase
	paymentpb "retarget/pkg/proto/payment"          // Импорт сгенерированного gRPC-кода

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	paymentUC usecase.PaymentUsecase // Добавляем зависимость
}

// NewPaymentServer — конструктор для инициализации сервера с зависимостями
func NewPaymentServer(paymentUC usecase.PaymentUsecase) *PaymentServer {
	return &PaymentServer{
		paymentUC: paymentUC,
	}
}

func RunGRPCServer(paymentUC usecase.PaymentUsecase) {
	lis, err := net.Listen("tcp", ":8054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(s, NewPaymentServer(paymentUC))

	log.Printf("gRPC server started on :8054")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *PaymentServer) RegUserActivity(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.PaymentResponse, error) {
	amount := entity.Decimal{}
	if err := amount.ParseFromString(req.GetAmount()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse amount: %v", err)
	}
	err := s.paymentUC.RegUserActivity(int(req.GetToUserId()), int(req.GetFromUserId()), amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	response := &paymentpb.PaymentResponse{
		TransactionId: strconv.Itoa(1),
		Status:        "success",
	}
	return response, nil
}
