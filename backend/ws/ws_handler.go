package ws

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{hub: h}
}

type CreateRoomRequest struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	DriverID   string `json:"driver_id"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:         req.ID,
		CustomerID: req.CustomerID,
		DriverID:   req.DriverID,
		Clients:    make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RoomInfoResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	DriverID   string `json:"driver_id"`
}

func (h *Handler) JoinRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	userID := c.Query("userId")
	phoneNumber := c.Query("phoneNumber")
	isCustomer := c.Query("isCustomer")

	resp, err := http.Get("http://localhost:6969/ws/room-info/" + roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var roomInfo RoomInfoResponse
	json.Unmarshal(responseData, &roomInfo)

	if (isCustomer == "true" && roomInfo.CustomerID != userID) || (isCustomer == "false" && roomInfo.DriverID != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "userID not correct - not authorized to join room"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client_id := userID
	if isCustomer == "true" {
		client_id = "cus_" + client_id
	} else {
		client_id = "drv_" + client_id
	}

	cl := &Client{
		Conn:           conn,
		Message:        make(chan *Message, 10),
		ID:             client_id,
		OriginalUserID: userID,
		RoomID:         roomID,
		PhoneNumber:    phoneNumber,
		IsCustomer:     isCustomer == "true",
	}

	// m := &Message{
	// 	Content:     "New user joined",
	// 	RoomID:      roomID,
	// 	UserID:      userID,
	// 	PhoneNumber: phoneNumber,
	// 	IsCustomer:  isCustomer == "true",
	// }

	// Register a new client through the register channel
	h.hub.Register <- cl
	// Broadcast that message
	// h.hub.BroadcastMessage <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *Handler) GetRoomInfo(c *gin.Context) {
	roomID := c.Param("roomId")
	if curr_room, ok := h.hub.Rooms[roomID]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
	} else {
		c.JSON(http.StatusOK, RoomInfoResponse{
			ID:         curr_room.ID,
			CustomerID: curr_room.CustomerID,
			DriverID:   curr_room.DriverID,
		})
	}
}

type ClientRes struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	IsCustomer  bool   `json:"is_customer"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:          c.ID,
			UserID:      c.OriginalUserID,
			PhoneNumber: c.PhoneNumber,
			IsCustomer:  c.IsCustomer,
		})
	}

	c.JSON(http.StatusOK, clients)
}
