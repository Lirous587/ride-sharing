package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse JSON data: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// validate
	if reqBody.UserID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal request", err), http.StatusInternalServerError)
		return
	}

	reader := bytes.NewBuffer(jsonBytes)

	// Call trip service
	resp, err := http.Post("http://trip-service:8083/trip/preview", "application/json", reader)
	if err != nil {
		http.Error(w, fmt.Sprintf("call trip server failed err: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var respBody any

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	response := contracts.APIResponse{
		Data: respBody,
	}

	util.WriteJson(w, resp.StatusCode, response)
}
