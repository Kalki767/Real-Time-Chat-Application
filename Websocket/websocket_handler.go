package websocket

import (
	"Real-Time-Chat-Application/domain"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development, configure appropriately for production
	},
}

// Client represents a websocket client connection
type Client struct {
	Conn     *websocket.Conn
	UserID   string
	ChatID   string
	SendChan chan []byte
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client] = true
			h.mutex.Unlock()

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.SendChan)
			}
			h.mutex.Unlock()

		case message := <-h.Broadcast:
			h.mutex.Lock()
			for client := range h.Clients {
				select {
				case client.SendChan <- message:
				default:
					close(client.SendChan)
					delete(h.Clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// HandleWebSocket upgrades HTTP connection to WebSocket and handles the connection
func HandleWebSocket(c *gin.Context, hub *Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	userID := c.Query("user_id")
	chatID := c.Query("chat_id")

	client := &Client{
		Conn:     conn,
		UserID:   userID,
		ChatID:   chatID,
		SendChan: make(chan []byte, 256),
	}

	hub.Register <- client

	// Start goroutines for reading and writing messages
	go client.writePump(hub)
	go client.readPump(hub)
}

func (c *Client) writePump(hub *Hub) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse the message
		var msg domain.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		// Broadcast the message to all clients in the same chat
		broadcastMsg := struct {
			ChatID  string         `json:"chat_id"`
			Message domain.Message `json:"message"`
		}{
			ChatID:  c.ChatID,
			Message: msg,
		}

		msgBytes, err := json.Marshal(broadcastMsg)
		if err != nil {
			log.Printf("error marshaling broadcast message: %v", err)
			continue
		}

		hub.Broadcast <- msgBytes
	}
}

// BroadcastToChat sends a message to all clients in a specific chat
func (h *Hub) BroadcastToChat(chatID string, message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.Clients {
		if client.ChatID == chatID {
			select {
			case client.SendChan <- message:
			default:
				close(client.SendChan)
				delete(h.Clients, client)
			}
		}
	}
}
