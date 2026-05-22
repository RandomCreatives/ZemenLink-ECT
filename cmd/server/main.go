package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"github.com/google/uuid"
	"zemenlink/internal/compliance"
	"zemenlink/internal/kernel"
	"zemenlink/internal/modules"
	"zemenlink/internal/tenant"
)

func main() {
	// 1. Setup Dependencies
	globalDBConn := os.Getenv("GLOBAL_DB_URL")
	if globalDBConn == "" {
		globalDBConn = "host=localhost user=postgres password=password dbname=zemenlink_global sslmode=disable"
	}

	globalDB, err := sqlx.Connect("postgres", globalDBConn)
	if err != nil {
		log.Fatalln("Failed to connect to global DB:", err)
	}

	masterKey := os.Getenv("MASTER_KEY")
	var kms tenant.KMSClient
	if masterKey != "" {
		kms, err = tenant.NewSymmetricKMS(masterKey)
		if err != nil {
			log.Fatalln("Failed to initialize KMS:", err)
		}
	}

	tenantManager := tenant.NewManager(globalDB, kms)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("default-secret-change-me")
	}

	// Load Modules
	modulesList, err := modules.LoadModules("./internal/modules")
	if err != nil {
		log.Println("Warning: Failed to load modules:", err)
	}
	moduleRegistry := modules.NewRegistry(modulesList)

	// Real-time Manager
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	rtManager := kernel.NewRTManager(redisURL)
	kernel.GlobalRTProvider = rtManager
	go rtManager.SubscribeToTenants(context.Background())

	// Compliance Engine & Handlers
	complianceEngine := &compliance.WorkflowEngine{}
	complianceHandlers := &compliance.Handlers{Engine: complianceEngine}

	// Crypto Services
	escrowCrypto, _ := kernel.NewEscrowCrypto(nil) // nil for PoC fallback
	kernel.GlobalEscrowCrypto = escrowCrypto

	// Messaging Pipeline
	msgPipeline := kernel.NewMessagePipeline()

	// 2. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(kernel.AuthMiddleware(jwtSecret, tenantManager, moduleRegistry))

		r.Get("/ws", rtManager.HandleWS)

		r.Post("/messages", func(w http.ResponseWriter, r *http.Request) {
			_, ok := kernel.GetTenantContext(r.Context())
			if !ok {
				http.Error(w, "Tenant context missing", http.StatusInternalServerError)
				return
			}

			var payload struct {
				ChatID       string `json:"chat_id"`
				Content      string `json:"content"`
				ContentType  string `json:"content_type"`
				RecipientID  string `json:"recipient_id"`
				EncryptedKey string `json:"encrypted_key"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "Invalid payload", http.StatusBadRequest)
				return
			}

			msgID := uuid.New().String()

			msg := &kernel.Message{
				ID:           msgID,
				Content:      payload.Content,
				RecipientID:  payload.RecipientID,
				EncryptedKey: payload.EncryptedKey,
				Metadata: map[string]interface{}{
					"chat_id":      payload.ChatID,
					"content_type": payload.ContentType,
				},
			}

			if err := msgPipeline.Execute(r.Context(), msg); err != nil {
				log.Printf("Pipeline error: %v", err)
				http.Error(w, "Failed to process message", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(msg)
		})

		r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
			tc, ok := kernel.GetTenantContext(r.Context())
			if !ok {
				http.Error(w, "Tenant context missing", http.StatusInternalServerError)
				return
			}

			// Example: Processing a message through the pipeline
			msg := &kernel.Message{Content: "Sample Enterprise Message"}
			if err := msgPipeline.Execute(r.Context(), msg); err != nil {
				http.Error(w, "Pipeline execution failed", http.StatusInternalServerError)
				return
			}

			var count int
			_ = tc.DB.Get(&count, "SELECT count(*) FROM messages") // Ignoring error for PoC
			fmt.Fprintf(w, "Tenant: %s. Compliance: %s. Modules: %v. Count: %d", tc.TenantID, tc.ComplianceLevel, tc.FeatureFlags, count)
		})

		// Compliance Routes
		complianceHandlers.RegisterRoutes(r)
	})

	// 3. Server Lifecycle
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("ZemenLink Backend starting on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done
	log.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}
