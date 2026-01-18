package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
	"ride-sharing/shared/util"
)

type HttpHandler struct {
	Service domain.TripService
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (h *HttpHandler) HandlePreviewTrip(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse JSON data: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	t, err := h.Service.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination)
	if err != nil {
		http.Error(w, fmt.Sprintf("create a trip failed,err:%v", err), http.StatusBadRequest)
		return
	}

	util.WriteJson(w, http.StatusOK, t)
}
