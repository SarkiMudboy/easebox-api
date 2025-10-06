package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SarkiMudboy/easebox-api/models"
	"github.com/gorilla/websocket"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var (
	ErrInvalidMessageType = &ValidationError{"invalid or missing type"}
	ErrMissingSessionID = &ValidationError{"missing session ID"}
	ErrMissingLocationData = &ValidationError{"location_update requires location data"}
	ErrorInvalidLatitude = &ValidationError{"latitude must be between -90 and 90"}
	ErrorInvalidLongitude = &ValidationError{"longitude must be between -180 and 180"}
	ErrorInvalidAccuracy = &ValidationError{"accuracy cannot be negative"}
)


var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict origins in production
		return true
	},
}


func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	defer conn.Close()

	log.Printf("New Websocket connection established: %v", r.RemoteAddr)

	conn.SetPongHandler(func (string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

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

	// Read loop
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket close error: %v", err)
			} else {
				log.Println("Websocket connection closed")
			}
			close(done)
			break
		}

		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		if err := handleMessage(message); err != nil {
			log.Printf("Error handling message: %v", err)
			sendError(conn, err.Error())
		}
	}
}

// Handles websocket messages from client
func handleMessage(message []byte) error {
	var wsMessage models.WebSocketMessage

	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return err
	}

	if err := validateMessage(&wsMessage); err != nil {
		return err
	}

	switch wsMessage.Type {
	case "start":
		handleStartTracking(&wsMessage)
	case "stop":
		handleStopTracking(&wsMessage)
	case "location_update":
		handleLocationUpdate(&wsMessage)
	default:
		log.Printf("Unknown message type: %v", wsMessage.Type)
	}

	return nil 
}


func validateMessage(message *models.WebSocketMessage) error {

	if message.Type == "" {
		return ErrInvalidMessageType
	}

	if message.SessionID == "" {
		return ErrMissingSessionID
	}

	if message.Type == "location_update" && message.Data == nil {
		return ErrMissingLocationData
	}
	if message.Data != nil {
		if message.Data.Latitude < -90 || message.Data.Latitude > 90 {
			return ErrorInvalidLatitude
		}
		if message.Data.Longitude < -180 || message.Data.Longitude > 180 {
			return ErrorInvalidLongitude
		}
		if message.Data.Accuracy < 0 {
			return ErrorInvalidAccuracy
		}
	}

	return nil
}

func handleStartTracking(msg *models.WebSocketMessage)  {
	log.Printf("[START] -> Session ID: %s, Started at : %v", msg.SessionID, time.Unix(0, *msg.State.StartTime*int64(time.Millisecond)))

}

func handleLocationUpdate(msg *models.WebSocketMessage) {
	log.Printf("[LOCATION UPDATE] -> Session ID: %s, Lat: %.6f, Lon: %.6f, Accuracy: %.2fm, Timestamp: %v",
		msg.SessionID,
		msg.Data.Latitude,
		msg.Data.Longitude,
		msg.Data.Accuracy,
		time.Unix(0, msg.Data.Timestamp*int64(time.Millisecond)),
	)

	if msg.Data.Speed > 0 {
		log.Printf("...Speed: %d m/s", msg.Data.Speed)
	}
	if msg.Data.Heading > 0 {
		log.Printf("...Heading: %d degrees", msg.Data.Heading)
	}

}

func handleStopTracking(msg *models.WebSocketMessage) {
	var duration time.Duration
	if msg.State.StartTime != nil && msg.State.LastUpdateTime != nil {
		startTime := time.Unix(0, *msg.State.StartTime * int64(time.Millisecond))
		endTime := time.Unix(0, *msg.State.LastUpdateTime * int64(time.Millisecond))
		duration = endTime.Sub(startTime)

		log.Printf("[STOP] Session ID: %v, Duration: %v", msg.SessionID, duration)

	}
}

func sendError(conn *websocket.Conn, errorMsg string) {
	errorResponse := map[string]string{
		"error": errorMsg,
	}

	if err := conn.WriteJSON(errorResponse); err != nil {
		log.Printf("Error sending error message to client: %v", err)
	}

}