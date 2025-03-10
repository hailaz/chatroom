package service

import (
	"chatroom/internal/consts"
	"chatroom/internal/dao"
	"chatroom/internal/model/entity"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

// WebSocketManager manages all WebSocket connections
type WebSocketManager struct {
	connections sync.Map // map[roomID]*sync.Map(map[userID]*Connection)
	upgrader    websocket.Upgrader
	userDao     *dao.UserDao
}

// Connection represents a WebSocket connection
type Connection struct {
	conn     *websocket.Conn
	user     *entity.User
	roomId   uint
	send     chan []byte
	manager  *WebSocketManager
	lastPing time.Time
}

// WebSocketMessage represents a message structure for WebSocket communication
type WebSocketMessage struct {
	Type      int         `json:"type"`
	Content   string      `json:"content,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
	UserId    uint        `json:"userId,omitempty"`
	Username  string      `json:"username,omitempty"`
	Nickname  string      `json:"nickname,omitempty"`
	Avatar    string      `json:"avatar,omitempty"`
}

var (
	wsManager *WebSocketManager
	wsOnce    sync.Once
)

// GetWebSocketManager returns the singleton instance of WebSocketManager
func GetWebSocketManager() *WebSocketManager {
	wsOnce.Do(func() {
		wsManager = &WebSocketManager{
			upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true // Allow all origins in development
				},
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			},
			userDao: dao.NewUserDao(),
		}

		// Start heartbeat checker
		go wsManager.heartbeatChecker()
	})
	return wsManager
}

// HandleWebSocket upgrades HTTP connection to WebSocket and handles the connection
func (m *WebSocketManager) HandleWebSocket(r *ghttp.Request, user *entity.User, roomId uint) {
	// Upgrade connection
	ws, err := m.upgrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Error(r.Context(), "WebSocket upgrade failed:", err)
		return
	}

	// Create new connection
	conn := &Connection{
		conn:     ws,
		user:     user,
		roomId:   roomId,
		send:     make(chan []byte, 256),
		manager:  m,
		lastPing: time.Now(),
	}

	// Update user status to online
	m.userDao.UpdateStatus(r.Context(), user.Id, consts.UserStatusOnline)

	// Store connection
	m.registerConnection(conn)

	// Start goroutines for reading and writing
	go conn.writePump()
	go conn.readPump()

	// Send current user list to all users in the room
	m.broadcastUserList(roomId)

	// Notify other users that new user joined
	m.broadcastToRoom(roomId, WebSocketMessage{
		Type:      consts.WsMsgTypeJoin,
		Content:   fmt.Sprintf("%s joined the room", user.Nickname),
		Timestamp: time.Now().Format(time.RFC3339),
		UserId:    user.Id,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
	})
}

// registerConnection registers a new WebSocket connection
func (m *WebSocketManager) registerConnection(conn *Connection) {
	// Get or create room connections map
	value, _ := m.connections.LoadOrStore(conn.roomId, &sync.Map{})
	roomConns := value.(*sync.Map)

	// Store connection in room
	roomConns.Store(conn.user.Id, conn)
}

// removeConnection removes a WebSocket connection
func (m *WebSocketManager) removeConnection(conn *Connection) {
	if value, ok := m.connections.Load(conn.roomId); ok {
		roomConns := value.(*sync.Map)
		roomConns.Delete(conn.user.Id)

		// Update user status to offline
		m.userDao.UpdateStatus(context.Background(), conn.user.Id, consts.UserStatusOffline)

		// If room is empty, remove it from connections map
		empty := true
		roomConns.Range(func(key, value interface{}) bool {
			empty = false
			return false
		})
		if empty {
			m.connections.Delete(conn.roomId)
		}

		// Broadcast user left message and updated user list
		m.broadcastToRoom(conn.roomId, WebSocketMessage{
			Type:      consts.WsMsgTypeLeave,
			Content:   fmt.Sprintf("%s离开了聊天室", conn.user.Nickname),
			Timestamp: time.Now().Format(time.RFC3339),
			UserId:    conn.user.Id,
			Username:  conn.user.Username,
			Nickname:  conn.user.Nickname,
			Avatar:    conn.user.Avatar,
		})

		// Broadcast updated user list
		m.broadcastUserList(conn.roomId)
	}
}

// broadcastToRoom broadcasts a message to all users in a room
func (m *WebSocketManager) broadcastToRoom(roomId uint, msg WebSocketMessage) {
	if value, ok := m.connections.Load(roomId); ok {
		roomConns := value.(*sync.Map)
		msgBytes, _ := json.Marshal(msg)

		roomConns.Range(func(key, value interface{}) bool {
			conn := value.(*Connection)
			select {
			case conn.send <- msgBytes:
			default:
				m.removeConnection(conn)
			}
			return true
		})
	}
}

// broadcastUserList sends the current user list to all users in a room
func (m *WebSocketManager) broadcastUserList(roomId uint) {
	if value, ok := m.connections.Load(roomId); ok {
		roomConns := value.(*sync.Map)

		// Build user list
		users := make([]map[string]interface{}, 0)
		roomConns.Range(func(key, value interface{}) bool {
			conn := value.(*Connection)
			users = append(users, map[string]interface{}{
				"id":       conn.user.Id,
				"username": conn.user.Username,
				"nickname": conn.user.Nickname,
				"avatar":   conn.user.Avatar,
				"status":   consts.UserStatusOnline,
			})
			return true
		})

		// Sort users by ID to maintain consistent order
		sort.Slice(users, func(i, j int) bool {
			return users[i]["id"].(uint) < users[j]["id"].(uint)
		})

		// Broadcast user list
		m.broadcastToRoom(roomId, WebSocketMessage{
			Type:      consts.WsMsgTypeUserList,
			Data:      users,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}

// heartbeatChecker periodically checks connection health
func (m *WebSocketManager) heartbeatChecker() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.connections.Range(func(roomId, value interface{}) bool {
			roomConns := value.(*sync.Map)
			roomConns.Range(func(userId, connValue interface{}) bool {
				conn := connValue.(*Connection)
				if time.Since(conn.lastPing) > 90*time.Second {
					// Connection is stale, close it
					conn.conn.Close()
					m.removeConnection(conn)
				}
				return true
			})
			return true
		})
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Connection) writePump() {
	pingTicker := time.NewTicker(30 * time.Second)
	defer func() {
		pingTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
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
		case <-pingTicker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Connection) readPump() {
	defer func() {
		c.manager.removeConnection(c)
		c.conn.Close()
		// Notify other users that user left
		c.manager.broadcastToRoom(c.roomId, WebSocketMessage{
			Type:      consts.WsMsgTypeLeave,
			Content:   fmt.Sprintf("%s left the room", c.user.Nickname),
			Timestamp: time.Now().Format(time.RFC3339),
			UserId:    c.user.Id,
			Username:  c.user.Username,
			Nickname:  c.user.Nickname,
			Avatar:    c.user.Avatar,
		})
	}()

	c.conn.SetReadLimit(1024 * 1024) // 1MB
	c.conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.lastPing = time.Now()
		c.conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				g.Log().Error(context.Background(), "WebSocket read error:", err)
			}
			break
		}

		// Update last ping time
		c.lastPing = time.Now()

		// Handle message
		var wsMsg WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		// Add user and timestamp information
		wsMsg.UserId = c.user.Id
		wsMsg.Username = c.user.Username
		wsMsg.Nickname = c.user.Nickname
		wsMsg.Avatar = c.user.Avatar
		wsMsg.Timestamp = time.Now().Format(time.RFC3339)

		// Store message in database if it's not a system message
		if wsMsg.Type != consts.MessageTypeSystem {
			messageDao := dao.NewMessageDao()
			messageDao.Create(context.Background(), &entity.Message{
				RoomId:  c.roomId,
				UserId:  c.user.Id,
				Content: wsMsg.Content,
				Type:    wsMsg.Type,
			})
		}

		// Broadcast message
		c.manager.broadcastToRoom(c.roomId, wsMsg)
	}
}
