package adv

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/pkg/entity"
	"time"

	"strconv"

	"github.com/mailru/easyjson"
)

func (c *AdvController) MetricsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	action := query.Get("action")
	bannerIDstr := query.Get("banner")
	slot := query.Get("slot")

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Format"))
		return
	}

	if err = c.advUsecase.WriteMetric(bannerID, slot, action); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got"))
}

func (c *AdvController) MyMetricsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fromStr := query.Get("from")
	toStr := query.Get("to")
	activity := query.Get("activity") // shown, click, ctr, уникальные для слотов: avg-action-price, revenue; баннера: expenses
	bannerIDstr := query.Get("banner")
	slotIDstr := query.Get("slot")

	if fromStr == "" || toStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "missing 'from' or 'to' parameters"))
		resp := entity.NewResponse(true, "missing 'from' or 'to' parameters")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	layout := "2006-01-02"

	fromTime, err := time.Parse(layout, fromStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "invalid 'to' format, use YYYY-MM-DD"))
		resp := entity.NewResponse(true, "invalid 'to' format, use YYYY-MM-DD")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	toTime, err := time.Parse(layout, toStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "invalid 'to' format, use YYYY-MM-DD"))
		resp := entity.NewResponse(true, "invalid 'to' format, use YYYY-MM-DD")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if fromTime.After(toTime) {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "'from' must be before or equal to 'to'"))
		resp := entity.NewResponse(true, "'from' must be before or equal to 'to'")
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

	if bannerIDstr != "" && slotIDstr == "" {
		bannerID, err := strconv.Atoi(bannerIDstr)
		if err != nil {
			http.Error(w, "Invalid banner ID", http.StatusBadRequest)
			return
		}
		var metrics interface{}
		if activity == "click" || activity == "shown" {
			metrics, err = c.advUsecase.GetBannerMetric(bannerID, activity, userID, fromTime, toTime)
		} else if activity == "ctr" {
			metrics, err = c.advUsecase.GetBannerCTR(bannerID, activity, userID, fromTime, toTime)
		} else if activity == "expenses" {
			metrics, err = c.advUsecase.GetBannerExpenses(bannerID, activity, userID, fromTime, toTime)
		} else {
			err = fmt.Errorf("unknown get parameters")
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//nolint:errcheck
			// json.NewEncoder(w).Encode(entity.NewResponse(true, "Bad GET parameters"))
			resp := entity.NewResponse(true, "Bad GET parameters")
			//nolint:errcheck
			easyjson.MarshalToWriter(&resp, w)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "metrics received", metrics))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	} else if bannerIDstr == "" && slotIDstr != "" {
		var metrics interface{}
		var err error
		err = nil
		if activity == "click" || activity == "shown" {
			metrics, err = c.advUsecase.GetSlotMetric(slotIDstr, activity, userID, fromTime, toTime)
		} else if activity == "ctr" {
			metrics, err = c.advUsecase.GetSlotCTR(slotIDstr, activity, userID, fromTime, toTime)
		} else if activity == "revenue" {
			metrics, err = c.advUsecase.GetSlotRevenue(slotIDstr, activity, userID, fromTime, toTime)
		} else if activity == "avg-show-price" {
			metrics, err = c.advUsecase.GetSlotAVGPrice(slotIDstr, activity, userID, fromTime, toTime)
		} else {
			err = fmt.Errorf("unknown get parameters")
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//nolint:errcheck
			// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
			resp := entity.NewResponse(true, err.Error())
			//nolint:errcheck
			easyjson.MarshalToWriter(&resp, w)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "metrics received", metrics))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(entity.NewResponse(true, "Bad GET parameters"))
	resp := entity.NewResponse(true, "Bad GET parameters")
	//nolint:errcheck
	easyjson.MarshalToWriter(&resp, w)

	// var metrics map[string]int
	// if bannerIDstr == "" && slotIDstr != "" {
	// 	metrics, err = c.advUsecase.GetSlotMetric(slotIDstr, activity, userID, fromTime, toTime)
	// }
	// if bannerIDstr != "" && slotIDstr == "" {
	// 	bannerID, convErr := strconv.Atoi(bannerIDstr)
	// 	if convErr != nil {
	// 		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
	// 		return
	// 	}
	// 	metrics, err = c.advUsecase.GetBannerMetric(bannerID, activity, userID, fromTime, toTime)
	// }
	// if err != nil {

	// }
	// w.WriteHeader(http.StatusOK)
	// err = json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "metrics received", metrics))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}
