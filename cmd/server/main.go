package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"zemenlink/internal/tenant"
)

func main() {
	// Initialize Global DB
	globalDBConn := os.Getenv("GLOBAL_DB_URL")
	if globalDBConn == "" {
		globalDBConn = "host=localhost user=postgres password=password dbname=zemenlink_global sslmode=disable"
	}

	globalDB, err := sqlx.Connect("postgres", globalDBConn)
	if err != nil {
		log.Fatalln("Failed to connect to global DB:", err)
	}

	tenantManager := tenant.NewManager(globalDB)

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract Tenant ID from Context (previously set by Auth Middleware)
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			http.Error(w, "Tenant ID required", http.StatusUnauthorized)
			return
		}

		// 2. Get the specific DB connection for this tenant
		db, err := tenantManager.GetDB(r.Context(), tenantID)
		if err != nil {
			http.Error(w, "Failed to connect to tenant database", http.StatusInternalServerError)
			return
		}

		// 3. Perform operations on the tenant database (Example)
		var count int
		err = db.Get(&count, "SELECT count(*) FROM messages")
		if err != nil {
			// In a real app, we'd handle this more gracefully
		}

		fmt.Fprintf(w, "Connected to database for tenant: %s. Message count: %d", tenantID, count)
	})

	log.Println("ZemenLink Backend starting on :8080...")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
