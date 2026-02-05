package messaging

import (
	pbd "ride-sharing/shared/proto/driver"
	pbt "ride-sharing/shared/proto/trip"
)

const (
	FindAvailableDriversQueue        = "find_available_drivers"
	DriverCmdTripRequestQueue        = "driver_cmd_trip_request"
	DriverTripResponseQueue          = "driver_trip_response"
	NotifyRiderNoDriversFoundQueue   = "notify_rider_no_drivers_found"
	NotifyDriverAssignQueue          = "notify_driver_assign"
	PaymentTripResponseQueue         = "payment_trip_response"
	NotifyPaymentSessionCreatedQueue = "notify_payment_session_created"
	NotifyPaymentSuccessQueue        = "payment_success"
)

type TripEventData struct {
	Trip    *pbt.Trip `json:"trip"`
	RiderID string    `json:"riderID"`
}

type DriverTripEventResponseData struct {
	Driver  *pbd.Driver `json:"driver"`
	TripID  string      `json:"tripID"`
	RiderID string      `json:"riderID"`
}

type PaymentEventSessionCreatedData struct {
	TripID    string  `json:"tripID"`
	SessionID string  `json:"sessionID"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type PaymentTripResponseData struct {
	TripID   string  `json:"tripID"`
	UserID   string  `json:"userID"`
	DriverID string  `json:"driverID"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PaymentStatusUpdateData struct {
	TripID   string `json:"tripID"`
	UserID   string `json:"userID"`
	DriverID string `json:"driverID"`
}
