package adv

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/pkg/entity"
	response "retarget/pkg/entity"
	"time"

	"strconv"
)

func (c *AdvController) MetricsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	action := query.Get("action")
	bannerIDstr := query.Get("banner")
	slot := query.Get("slot")

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Format"))
		return
	}

	if err = c.advUsecase.WriteMetric(bannerID, slot, action); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got"))
}

func (c *AdvController) MyMetricsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fromStr := query.Get("from")
	toStr := query.Get("to")
	activity := query.Get("activity")
	bannerIDstr := query.Get("banner")
	slotIDstr := query.Get("slot")

	if fromStr == "" || toStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "missing 'from' or 'to' parameters"))
		return
	}

	// Формат даты
	layout := "2006-01-02"

	// Парсинг
	fromTime, err := time.Parse(layout, fromStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "invalid 'to' format, use YYYY-MM-DD"))
		return
	}
	toTime, err := time.Parse(layout, toStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "invalid 'to' format, use YYYY-MM-DD"))
		return
	}

	// if fromTime.After(toTime) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(entity.NewResponse(true, "'from' must be before or equal to 'to'"))
	// 	return
	// }

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	if bannerIDstr == "" && slotIDstr != "" {
		metrics, err := c.advUsecase.GetSlotMetric(slotIDstr, activity, userID, fromTime, toTime)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "metrics received", metrics))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	if bannerIDstr != "" && slotIDstr == "" {
		bannerID, err := strconv.Atoi(bannerIDstr)
		metrics, err := c.advUsecase.GetBannerMetric(bannerID, activity, userID, fromTime, toTime)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "metrics received", metrics))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
}
