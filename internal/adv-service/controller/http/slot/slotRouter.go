package slot

import (
	"net/http"
	usecaseSlot "retarget/internal/adv-service/usecase/slot"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

type SlotController struct {
	slotUsecase *usecaseSlot.SlotUsecase
}

func NewSlotController(slotUsecase *usecaseSlot.SlotUsecase) *SlotController {
	return &SlotController{slotUsecase: slotUsecase}
}

func SetupSlotRoutes(authenticator *authenticate.Authenticator, slotUsecase *usecaseSlot.SlotUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	slotController := NewSlotController(slotUsecase)

	muxRouter.Use(authenticate.AuthMiddleware(authenticator))

	muxRouter.Handle("/api/v1/slot/create", http.HandlerFunc(slotController.CreateSlotHandler)).Methods("POST")
	muxRouter.Handle("/api/v1/slot/edit", http.HandlerFunc(slotController.EditSlotHandler)).Methods("PUT")
	muxRouter.Handle("/api/v1/slot/delete", http.HandlerFunc(slotController.DeleteSlotHandler)).Methods("DELETE")
	muxRouter.Handle("/api/v1/slot/my", http.HandlerFunc(slotController.GetUserSlotsHandler)).Methods("GET")
	muxRouter.Handle("/api/v1/slot/formats", http.HandlerFunc(slotController.GetFormatsHandler)).Methods("GET")

	return muxRouter
}
