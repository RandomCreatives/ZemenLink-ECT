package tenant

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Tenant represents a client organization
type Tenant struct {
	ID               string `db:"id"`
	DBConnectionString string `db:"db_connection_string"`
	ComplianceLevel   string `db:"compliance_level"` // e.g., "high", "standard"
}

// Manager handles connections to tenant-specific databases
type Manager struct {
	globalDB *sqlx.DB
	conns    map[string]*sqlx.DB
	mu       sync.RWMutex
}

func NewManager(globalDB *sqlx.DB) *Manager {
	return &Manager{
		globalDB: globalDB,
		conns:    make(map[string]*sqlx.DB),
	}
}

// GetDB returns a database connection for the given tenant ID
func (m *Manager) GetDB(ctx context.Context, tenantID string) (*sqlx.DB, error) {
	m.mu.RLock()
	db, ok := m.conns[tenantID]
	m.mu.RUnlock()
	if ok {
		return db, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring lock
	if db, ok = m.conns[tenantID]; ok {
		return db, nil
	}

	// Lookup connection string from Global DB
	var tenant Tenant
	err := m.globalDB.GetContext(ctx, &tenant, "SELECT id, db_connection_string FROM tenants WHERE id = $1", tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant: %w", err)
	}

	// Decrypt the connection string (Implementation would depend on KMS choice)
	decryptedConn, err := m.decryptConnectionString(tenant.DBConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt connection string: %w", err)
	}

	// Initialize new connection
	newDB, err := sqlx.Open("postgres", decryptedConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open tenant db: %w", err)
	}

	// Optional: Configure pooling
	newDB.SetMaxOpenConns(20)
	newDB.SetMaxIdleConns(5)

	m.conns[tenantID] = newDB
	return newDB, nil
}

func (m *Manager) decryptConnectionString(encrypted string) (string, error) {
	// TODO: Integrate with KMS (e.g. AWS KMS, HashiCorp Vault)
	// For now, returning as-is for the PoC
	return encrypted, nil
}
