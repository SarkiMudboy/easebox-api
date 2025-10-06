package models

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp int64   `json:"timestamp"`
	Speed     int64   `json:"speed"`
	Heading   int64   `json:"heading"`
}

type TrackingState struct {
	IsTracking     bool   `json:"isTracking"`
	SessionID      string `json:"sessionId"`
	StartTime      *int64 `json:"startTime"`
	LastUpdateTime *int64 `json:"lastUpdateTime"`
}

type WebSocketMessage struct {
	Type      string        `json:"type"`
	SessionID string        `json:"sessionId"`
	Data      *LocationData `json:"data"`
	State     TrackingState `json:"state"`
}