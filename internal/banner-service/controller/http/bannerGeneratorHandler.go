package controller

import (
	"fmt"
	"github.com/mailru/easyjson"
	"net/http"
	response "retarget/pkg/entity"
)

func (h *BannerController) GenerateDescription(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	title := query.Get("title")

	description, err := h.BannerUsecase.GenerateBannerDescription(title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := response.NewResponse(false, description)
	//nolint:errcheck
	easyjson.MarshalToWriter(&resp, w)
}

func (h *BannerController) GenerateImage(w http.ResponseWriter, r *http.Request) {
	//requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	//userCtx, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	//if !ok {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	easyjson.MarshalToWriter(response.NewResponse(true, "Ошибка аутентификации"), w)
	//	return
	//}

	query := r.URL.Query()
	title := query.Get("title")
	// userID := userCtx.UserID

	//id, err := strconv.Atoi(mux.Vars(r)["banner_id"])
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	easyjson.MarshalToWriter(response.NewResponse(true, "Некорректный ID баннера"), w)
	//	return
	//}

	imgBytes, err := h.BannerUsecase.GenerateBannerImage(title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		easyjson.MarshalToWriter(response.NewResponse(true, err.Error()), w)
		return
	}

	contentType := http.DetectContentType(imgBytes)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imgBytes)))
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=banner-%d.png"))
	w.WriteHeader(http.StatusOK)

	w.Write(imgBytes)
}
