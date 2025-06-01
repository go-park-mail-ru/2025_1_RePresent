package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	model "retarget/internal/pay-service/easyjsonModels"
	"retarget/pkg/entity"
	response "retarget/pkg/entity"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/google/uuid"
)

func (h *PaymentController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Cookie"))
		resp := entity.NewResponse(true, "Invalid Cookie")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		resp := entity.NewResponse(true, "Error of authenticator")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	userID := userSession.UserID

	balance, err := h.PaymentUsecase.GetBalanceByUserId(userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(
		// 	true,
		// 	"Error fetching balance: "+err.Error(),
		// ))
		resp := entity.NewResponse(true, "Error fetching balance: "+err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseData := map[string]interface{}{
		"balance": balance,
	}

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(
		// 	true,
		// 	"Error encoding response: "+err.Error(),
		// ))

		resp := entity.NewResponse(true, "Error encoding response: "+err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
	}
}

func (h *PaymentController) TopUpAccount(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Cookie"))
		resp := entity.NewResponse(true, "Invalid Cookie")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		resp := entity.NewResponse(true, "Error of authenticator")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	var req model.TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Request Body"))

		resp := entity.NewResponse(true, "Invalid Request Body")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Amount"))
		resp := entity.NewResponse(true, "Invalid Amount")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if err = h.PaymentUsecase.TopUpBalance(userID, req.Amount, requestID); err != nil {
		// handleTopUpError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	transactionID := uuid.New().String()

	responseData := model.TransactionResponse{
		TransactionID: transactionID,
		Status:        "completed",
		NextAction:    "redirect_to_payment_gateway",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	//nolint:errcheck
	easyjson.MarshalToWriter(&responseData, w)

}

func (c *PaymentController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Value          string `json:"value"`
		Currency       string `json:"currency"`
		ReturnURL      string `json:"return_url"`
		Description    string `json:"description"`
		IdempotenceKey string `json:"idempotence_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Cookie"))
		return
	}

	// userID := r.Context().Value("user_id").(int)

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		return
	}
	userID := userSession.UserID

	confirmationURL, err := c.PaymentUsecase.CreateYooMoneyPayment(
		userID,
		req.Value,
		req.Currency,
		req.ReturnURL,
		req.Description,
		req.IdempotenceKey,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//nolint:errcheck
	json.NewEncoder(w).Encode(map[string]string{"confirmation_url": confirmationURL})
}

/* Хз мб это нужно
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
*/

func (h *PaymentController) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	// transactionID := r.URL.Query().Get("transactionId")

	vars := mux.Vars(r)
	transactionID := vars["transactionid"]

	tx, err := h.PaymentUsecase.GetTransactionByID(transactionID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(tx)
}

func (h *PaymentController) RegUserActivity(w http.ResponseWriter, r *http.Request) {

}

func (c *PaymentController) WithdrawFunds(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)

	// теперь принимаем return_url
	var req struct {
		Amount      float64 `json:"amount"`
		ReturnURL   string  `json:"return_url"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Amount < 10.0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Minimum withdrawal amount is 10.00"))
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		return
	}
	userID := userSession.UserID

	idemKey := fmt.Sprintf("payout_%s_%d_%s", requestID, userID, uuid.New().String())

	// redirect-flow
	redirectURL, err := c.PaymentUsecase.CreateYooMoneyPayoutRedirect(
		userID, req.Amount, req.Description, req.ReturnURL, idemKey,
	)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "insufficient funds") {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"confirmation_url": redirectURL,
	})
}

func (c *PaymentController) WithdrawFundsRedirect(w http.ResponseWriter, r *http.Request) {
	// _ := r.Context().Value(response.СtxKeyRequestID{}).(string)

	var req struct {
		Amount      float64 `json:"amount"`
		ReturnURL   string  `json:"return_url"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Amount < 10.0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Minimum withdrawal amount is 10.00"))
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		return
	}
	userID := userSession.UserID

	// Генерируем более короткий идемпотентный ключ
	// Используем только первые 8 символов UUID и ID пользователя
	shortUUID := strings.ReplaceAll(uuid.New().String()[:8], "-", "")
	idemKey := fmt.Sprintf("p%d%s", userID, shortUUID)

	// redirect-flow
	redirectURL, err := c.PaymentUsecase.CreateYooMoneyPayoutRedirect(
		userID, req.Amount, req.Description, req.ReturnURL, idemKey,
	)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "insufficient funds") {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"confirmation_url": redirectURL,
	})
}
