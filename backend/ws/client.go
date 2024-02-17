package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn        *websocket.Conn
	Message     chan *Message
	ID          string `json:"id"`
	RoomID      string `json:"room_id"`
	PhoneNumber string `json:"phone_number"`
	IsCustomer  bool   `json:"is_customer"`
}

type Message struct {
	Content     string `json:"content"`
	RoomID      string `json:"room_id"`
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	IsCustomer  bool   `json:"is_customer"`
}

type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:     string(m),
			RoomID:      c.RoomID,
			UserID:      c.ID,
			PhoneNumber: c.PhoneNumber,
			IsCustomer:  c.IsCustomer,
		}

		hub.Broadcast <- msg
	}
}
