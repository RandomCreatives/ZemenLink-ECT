package kernel

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, check against allowed domains
	},
}

func (m *RTManager) HandleWS(w http.ResponseWriter, r *http.Request) {
	tc, ok := GetTenantContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := tc.UserID
	if userID == "" {
		// Strictly enforce identity from JWT context
		http.Error(w, "Unauthorized: user_id missing from context", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}

	m.AddConn(tc.TenantID, userID, conn)
	defer m.RemoveConn(tc.TenantID, userID)
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket read error: %v", err)
			}
			break
		}
	}
}
