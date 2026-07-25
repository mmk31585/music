package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// AuthorizeSub provides channel-level authorization — return true to allow subscription.
// If nil, all subscriptions are allowed.
type AuthorizeSub func(userID uuid.UUID, channel string) bool

type Message struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan []byte
	mu     sync.Mutex
}

type Subscription struct {
	Channel string
	Client  *Client
}

type Hub struct {
	mu            sync.RWMutex
	clients       map[uuid.UUID]*Client
	Authorize     AuthorizeSub
	subscriptions map[string]map[*Client]bool
	broadcast     chan BroadcastMessage
	register      chan *Client
	unregister    chan *Client
	logger        *zap.Logger
	upgrader      websocket.Upgrader
}

type BroadcastMessage struct {
	Channel string
	Data    []byte
	Exclude uuid.UUID
}

// NewHub creates a WebSocket hub.
// If allowedOrigins is empty, only same-origin requests are permitted (Origin header must be empty).
func NewHub(logger *zap.Logger, allowedOrigins []string) *Hub {
	originChecker := buildOriginChecker(allowedOrigins)

	return &Hub{
		clients:       make(map[uuid.UUID]*Client),
		subscriptions: make(map[string]map[*Client]bool),
		broadcast:     make(chan BroadcastMessage, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		logger:        logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:    4096,
			WriteBufferSize:   4096,
			EnableCompression: true,
			CheckOrigin:       originChecker,
		},
	}
}

// buildOriginChecker returns a CheckOrigin func that validates against the allowed origins list.
// It supports:
//   - Exact match against each allowed origin
//   - Same-origin requests (empty Origin header) are allowed when the list is non-empty
//   - If allowedOrigins is empty, only same-origin (no Origin header) is permitted
func buildOriginChecker(allowedOrigins []string) func(r *http.Request) bool {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o != "" {
			allowed[o] = struct{}{}
		}
	}

	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		// No Origin header = same-origin request (browser doesn't send it for same-origin)
		if origin == "" {
			return true
		}

		// Empty allowed list means only same-origin
		if len(allowed) == 0 {
			return false
		}

		_, ok := allowed[origin]
		return ok
	}
}

func (h *Hub) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			h.shutdown()
			return ctx.Err()

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.logger.Info("ws client connected", zap.String("user_id", client.UserID.String()))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			for channel, subs := range h.subscriptions {
				delete(subs, client)
				if len(subs) == 0 {
					delete(h.subscriptions, channel)
				}
			}
			h.mu.Unlock()
			h.logger.Info("ws client disconnected", zap.String("user_id", client.UserID.String()))

		case msg := <-h.broadcast:
			h.mu.RLock()
			subs, ok := h.subscriptions[msg.Channel]
			if !ok {
				h.mu.RUnlock()
				continue
			}
			for client := range subs {
				if client.UserID == msg.Exclude {
					continue
				}
				select {
				case client.Send <- msg.Data:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Subscribe(client *Client, channel string) bool {
	if h.Authorize != nil && !h.Authorize(client.UserID, channel) {
		h.logger.Warn("ws subscribe denied",
			zap.String("user_id", client.UserID.String()),
			zap.String("channel", channel))
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscriptions[channel] == nil {
		h.subscriptions[channel] = make(map[*Client]bool)
	}
	h.subscriptions[channel][client] = true
	return true
}

func (h *Hub) Unsubscribe(client *Client, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if subs, ok := h.subscriptions[channel]; ok {
		delete(subs, client)
		if len(subs) == 0 {
			delete(h.subscriptions, channel)
		}
	}
}

func (h *Hub) BroadcastToChannel(channel string, message Message, excludeUserID uuid.UUID) {
	data, _ := json.Marshal(message)
	h.broadcast <- BroadcastMessage{
		Channel: channel,
		Data:    data,
		Exclude: excludeUserID,
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	data, _ := json.Marshal(message)
	select {
	case client.Send <- data:
	default:
	}
}

func (h *Hub) OnlineUsers() []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := make([]uuid.UUID, 0, len(h.clients))
	for id := range h.clients {
		users = append(users, id)
	}
	return users
}

func (h *Hub) IsOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) shutdown() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients {
		close(client.Send)
	}
	h.logger.Info("ws hub shut down, all clients notified")
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("ws upgrade failed", zap.Error(err))
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Hub:    h,
		Send:   make(chan []byte, 256),
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(4096)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "subscribe":
			if channel, ok := msg.Payload.(string); ok {
				if !c.Hub.Subscribe(c, channel) {
					c.Send <- []byte(`{"type":"error","payload":"subscription denied"}`)
				}
			} else {
				c.Send <- []byte(`{"type":"error","payload":"invalid channel"}`)
			}
		case "unsubscribe":
			if channel, ok := msg.Payload.(string); ok {
				c.Hub.Unsubscribe(c, channel)
			}
		case "ping":
			c.Send <- []byte(`{"type":"pong"}`)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "shutdown"))
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "hub shutting down"))
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
