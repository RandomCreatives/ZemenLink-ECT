package kernel

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// RTManager handles real-time connections and pub/sub
type RTManager struct {
	redisClient *redis.Client
	conns      map[string]map[string]*websocket.Conn // tenantID -> userID -> conn
	mu         sync.RWMutex
}

func NewRTManager(redisURL string) *RTManager {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: Failed to parse Redis URL: %v. Real-time features may be limited.", err)
		return &RTManager{conns: make(map[string]map[string]*websocket.Conn)}
	}

	client := redis.NewClient(opt)
	return &RTManager{
		redisClient: client,
		conns:      make(map[string]map[string]*websocket.Conn),
	}
}

func (m *RTManager) AddConn(tenantID, userID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.conns[tenantID]; !ok {
		m.conns[tenantID] = make(map[string]*websocket.Conn)
	}
	m.conns[tenantID][userID] = conn
}

func (m *RTManager) RemoveConn(tenantID, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if userConns, ok := m.conns[tenantID]; ok {
		delete(userConns, userID)
		if len(userConns) == 0 {
			delete(m.conns, tenantID)
		}
	}
}

// BroadcastToTenant publishes a message to a Redis channel for the tenant
func (m *RTManager) BroadcastToTenant(ctx context.Context, tenantID string, payload interface{}) error {
	if m.redisClient == nil {
		// Fallback to local broadcast if Redis is missing
		m.localBroadcast(tenantID, payload)
		return nil
	}
	return m.redisClient.Publish(ctx, "tenant:"+tenantID, payload).Err()
}

func (m *RTManager) localBroadcast(tenantID string, payload interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userConns, ok := m.conns[tenantID]
	if !ok {
		return
	}

	for userID, conn := range userConns {
		// Run in goroutine to prevent a slow consumer from blocking the broadcast
		go func(uid string, c *websocket.Conn) {
			if err := c.WriteJSON(payload); err != nil {
				log.Printf("Error writing to websocket (user: %s): %v", uid, err)
			}
		}(userID, conn)
	}
}

// SubscribeToTenants starts listening for Redis messages for active tenants
func (m *RTManager) SubscribeToTenants(ctx context.Context) {
	if m.redisClient == nil {
		return
	}

	pubsub := m.redisClient.PSubscribe(ctx, "tenant:*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		// msg.Channel is "tenant:<id>"
		var tenantID string
		fmt.Sscanf(msg.Channel, "tenant:%s", &tenantID)

		m.localBroadcast(tenantID, msg.Payload)
	}
}
