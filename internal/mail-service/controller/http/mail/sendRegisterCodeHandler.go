package mail

import (
	"encoding/json"
	"net/http"
	entityMail "retarget/internal/mail-service/entity/mail"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"
	"strings"
)

type RegisterCodeRequest struct {
	Email string `json:"email" validate:"email,required"`
	Code  string `json:"code" validate:"required,len=6"`
}

func (c *MailController) SendRegisterCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	var registerCodeRequest RegisterCodeRequest
	err := json.NewDecoder(r.Body).Decode(&registerCodeRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid request body"))
		return
	}

	errorMessages, err := validator.ValidateStruct(registerCodeRequest)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, errorMessages))
		return
	}

	err = c.mailUsecase.SendCodeMail(entityMail.REGISTER, registerCodeRequest.Email, registerCodeRequest.Code)
	if err != nil {
		if strings.HasPrefix(err.Error(), "5") {
			w.WriteHeader(http.StatusBadRequest)
			//nolint:errcheck
			json.NewEncoder(w).Encode(entity.NewResponse(true, "Такой почты не существует"))
			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Ошибка, повторите отправку позже"))
		return
	}

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Sent"))
}
