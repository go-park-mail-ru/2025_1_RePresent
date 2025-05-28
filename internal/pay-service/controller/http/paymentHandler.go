package payment

import (
	"encoding/json"
	"net/http"
	model "retarget/internal/pay-service/easyjsonModels"
	"retarget/pkg/entity"
	response "retarget/pkg/entity"

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
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		resp := entity.NewResponse(true, "Error of authenticator")
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
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		resp := entity.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	var req model.TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Request Body"))

		resp := entity.NewResponse(true, "Invalid Request Body")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Amount"))
		resp := entity.NewResponse(true, "Invalid Amount")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if err = h.PaymentUsecase.TopUpBalance(userID, req.Amount, requestID); err != nil {
		// handleTopUpError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
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
