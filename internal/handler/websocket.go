package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SarkiMudboy/easebox-api/internal/domain"
	"github.com/SarkiMudboy/easebox-api/internal/service"
	"github.com/gorilla/websocket"
)

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp int64   `json:"timestamp"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
}

type TrackingState struct {
	IsTracking     bool   `json:"isTracking"`
	SessionID      string `json:"sessionId"`
	DeliveryID 	   string `json:"deliveryId"`
	StartTime      *int64 `json:"startTime"`
	LastUpdateTime *int64 `json:"lastUpdateTime"`
}

type WebSocketMessage struct {
	Type      string        `json:"type"`
	SessionID string        `json:"sessionId"`
	Data      *LocationData `json:"data"`
	State     TrackingState `json:"state"`
}

type WebSocketHandler struct {
	locationService *service.LocationService
	upgrader websocket.Upgrader
}

func NewWebSocketHandler(locationService *service.LocationService) *WebSocketHandler {
	return &WebSocketHandler{
		locationService: locationService,
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	log.Printf("New Websocket connection established: %v", r.RemoteAddr)

	defer conn.Close()


	conn.SetPongHandler(func (string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// ping op to maintain connection
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Channel for recieving stop signal: stop websocket
	done := make(chan bool)

	go func() {
		for {
			select {
			case <- ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10 * time.Second)); err != nil {
					log.Printf("Error sending ping: %v", err)
					return
				}
			case <- done:
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			close(done)
			break
		}

		var msg WebSocketMessage
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message content: %v", err)
			continue
		}

		ctx := r.Context()


		switch msg.Type {
		case "start":
			err = h.locationService.StartTracking(ctx, msg.State.SessionID, msg.State.DeliveryID)
			log.Printf("[START] -> Session ID: %s, Started at : %v", msg.SessionID, time.UnixMilli(*msg.State.StartTime))
		case "location_update":
			loc := h.MessageToLocation(&msg)

			if loc == nil {
				log.Printf("missing location_update data")
				continue
			}

			err = h.locationService.RecordLocation(ctx, loc)

			log.Printf("[LOCATION UPDATE] -> Session ID: %s, Lat: %.6f, Lon: %.6f, Accuracy: %.2fm, Timestamp: %v",
				msg.SessionID,
				msg.Data.Latitude,
				msg.Data.Longitude,
				msg.Data.Accuracy,
				time.UnixMilli(msg.Data.Timestamp),
			)

		case "stop":
			err = h.locationService.StopTracking(ctx, msg.SessionID)

			var duration time.Duration
			if msg.State.StartTime != nil && msg.State.LastUpdateTime != nil {
				startTime := time.UnixMilli(*msg.State.StartTime)
				endTime := time.UnixMilli(*msg.State.LastUpdateTime)
				duration = endTime.Sub(startTime)

				log.Printf("[STOP] Session ID: %v, Duration: %v", msg.SessionID, duration)

			}
		}

		if err != nil {
			log.Println(msg.Data, msg.State)
			log.Printf("Service error: %v", err)
		}

	}
}

func (h *WebSocketHandler) MessageToLocation (message *WebSocketMessage) *domain.LocationUpdate {

	if message.Data == nil {
		return nil
	}

	return &domain.LocationUpdate{
		SessionID: message.SessionID,
		DeliveryID: message.State.DeliveryID,
		Latitude: message.Data.Latitude,
		Longitude: message.Data.Longitude,
		Accuracy: message.Data.Accuracy,
		Speed: &message.Data.Speed,
		Heading: &message.Data.Heading,
		RecordedAt: time.UnixMilli(message.Data.Timestamp),
	}

}

