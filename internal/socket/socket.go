package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type EventHandler func(*Client, []byte)

type Hub struct {
	mu       sync.RWMutex
	clients  map[*Client]bool
	handlers map[string][]EventHandler
}

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[*Client]bool),
		handlers: make(map[string][]EventHandler),
	}
}

func (h *Hub) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("socket upgrade failed: %v", err)
		return
	}

	clientID := uuid.NewString()
	client := &Client{
		id:   clientID,
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	log.Printf("socket connected: user connected with id %s", clientID)

	if err := conn.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, "connected:%s", clientID)); err != nil {
		_ = conn.Close()
		return
	}

	go client.writePump()
	client.readPump()
}

func (h *Hub) On(event string, handler EventHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[event] = append(h.handlers[event], handler)
}

func (h *Hub) Emit(event string, payload any) {
	message, err := json.Marshal(map[string]any{
		"event": event,
		"data":  payload,
	})
	if err != nil {
		log.Printf("socket emit marshal failed: %v", err)
		return
	}

	h.Broadcast(string(message))
}

func (h *Hub) Broadcast(message string) {
	payload := []byte(message)

	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.send <- payload:
		default:
		}
	}
}

func (h *Hub) handleMessage(client *Client, message []byte) {
	var packet struct {
		Event string          `json:"event"`
		Data  json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(message, &packet); err != nil || packet.Event == "" {
		h.Broadcast(string(message))
		return
	}

	h.mu.RLock()
	handlers := append([]EventHandler(nil), h.handlers[packet.Event]...)
	h.mu.RUnlock()

	for _, handler := range handlers {
		handler(client, packet.Data)
	}
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	delete(h.clients, client)
	h.mu.Unlock()
	log.Printf("socket disconnected: user disconnected with id %s", client.id)
}

func (c *Client) readPump() {
	defer func() {
		c.hub.removeClient(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(1024 * 1024)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("socket read error: %v", err)
			}
			return
		}

		c.hub.handleMessage(c, message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
				return
			}
		}
	}
}
