package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketServer struct {
	mu          sync.RWMutex
	connections map[string]*WebSocketConn
	handler    WebSocketHandler
	server     *http.Server
	mux        *http.ServeMux
}

type WebSocketConn struct {
	ID       string
	conn     *websocket.Conn
	send     chan []byte
	closed   bool
	server   *WebSocketServer
}

type WebSocketHandler interface {
	HandleConnect(conn *WebSocketConn)
	HandleMessage(conn *WebSocketConn, message []byte)
	HandleDisconnect(conn *WebSocketConn)
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WSResponse struct {
	Type    string `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewWebSocketServer(addr string, handler WebSocketHandler) *WebSocketServer {
	ws := &WebSocketServer{
		connections: make(map[string]*WebSocketConn),
		handler:    handler,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.handleWebSocket)

	ws.mux = mux
	ws.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return ws
}

func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	wsConn := &WebSocketConn{
		conn:   conn,
		send:   make(chan []byte, 256),
		server: s,
	}

	s.mu.Lock()
	s.connections[wsConn.ID] = wsConn
	s.mu.Unlock()

	if s.handler != nil {
		s.handler.HandleConnect(wsConn)
	}

	go wsConn.writePump()
	go wsConn.readPump()
}

func (c *WebSocketConn) readPump() {
	defer func() {
		c.Close()
		if c.server.handler != nil {
			c.server.handler.HandleDisconnect(c)
		}
	}()

	c.conn.SetReadLimit(65536)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		if c.server.handler != nil {
			c.server.handler.HandleMessage(c, message)
		}
	}
}

func (c *WebSocketConn) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WebSocketConn) Send(messageType string, payload interface{}) error {
	resp := WSResponse{
		Type:    messageType,
		Payload: payload,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("send buffer full")
	}
}

func (c *WebSocketConn) SendError(err error) error {
	resp := WSResponse{
		Type:  "error",
		Error: err.Error(),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("send buffer full")
	}
}

func (c *WebSocketConn) Close() error {
	c.server.mu.Lock()
	if c.closed {
		c.server.mu.Unlock()
		return nil
	}
	c.closed = true
	delete(c.server.connections, c.ID)
	c.server.mu.Unlock()

	close(c.send)
	return c.conn.Close()
}

func (s *WebSocketServer) Broadcast(messageType string, payload interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.Marshal(WSResponse{
		Type:    messageType,
		Payload: payload,
	})
	if err != nil {
		return err
	}

	for _, conn := range s.connections {
		select {
		case conn.send <- data:
		default:
		}
	}

	return nil
}

func (s *WebSocketServer) SendTo(connID string, messageType string, payload interface{}) error {
	s.mu.RLock()
	conn, ok := s.connections[connID]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("connection not found: %s", connID)
	}

	return conn.Send(messageType, payload)
}

func (s *WebSocketServer) ConnectionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.connections)
}

func (s *WebSocketServer) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.Shutdown()
	}()

	return s.server.ListenAndServe()
}

func (s *WebSocketServer) Shutdown() error {
	s.mu.Lock()
	for _, conn := range s.connections {
		conn.Close()
	}
	s.mu.Unlock()

	return s.server.Close()
}

func (s *WebSocketServer) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *WebSocketServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}
