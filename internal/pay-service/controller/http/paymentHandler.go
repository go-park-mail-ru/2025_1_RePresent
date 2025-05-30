package payment

import (
	"encoding/json"
	"net/http"
	"retarget/internal/pay-service/repo"
	entity "retarget/pkg/entity"
	response "retarget/pkg/entity"

	"github.com/google/uuid"
)

type TransactionResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	NextAction    string `json:"nextAction"`
}

type TopUpRequest struct {
	Amount int64 `json:"amount"`
}

func (h *PaymentController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Cookie"))
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		return
	}
	userID := userSession.UserID

	balance, err := h.PaymentUsecase.GetBalanceByUserId(userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.NewResponse(
			true,
			"Error fetching balance: "+err.Error(),
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseData := map[string]interface{}{
		"balance": balance,
	}

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(
			true,
			"Error encoding response: "+err.Error(),
		))
	}
}

func (h *PaymentController) TopUpAccount(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Cookie"))
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	var req TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Request Body"))
		return
	}

	if req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Amount"))
		return
	}

	if err = h.PaymentUsecase.TopUpBalance(userID, req.Amount, requestID); err != nil {
		// handleTopUpError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	transactionID := uuid.New().String()

	responseData := TransactionResponse{
		TransactionID: transactionID,
		Status:        "completed",
		NextAction:    "redirect_to_payment_gateway",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(responseData)

}

func handleTopUpError(w http.ResponseWriter, err error) {
	switch err {
	case repo.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "User  not found"))
	case repo.ErrInvalidAmount:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid amount"))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(
			true,
			"Internal Server Error: "+err.Error(),
		))
	}
}

func (h *PaymentController) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	transactionID := r.URL.Query().Get("transactionId")

	tx, err := h.PaymentUsecase.GetTransactionByID(transactionID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tx)
}

func (h *PaymentController) RegUserActivity(w http.ResponseWriter, r *http.Request) {

}
