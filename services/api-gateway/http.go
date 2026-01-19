package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade websocket failed: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("No user ID provided")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Printf("No package slug provided")
		return
	}

	type Driver struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:             "123",
			Name:           "Lirous",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABC123",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Receive message: %s", msg)
	}
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("No user ID provided")
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Receive message: %s", msg)
	}
}
