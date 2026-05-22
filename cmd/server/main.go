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
			tc, ok := kernel.GetTenantContext(r.Context())
			if !ok {
				http.Error(w, "Tenant context missing", http.StatusInternalServerError)
				return
			}

			// In a real app, parse from JSON body
			msg := &kernel.Message{
				ID:      "msg-" + time.Now().String(),
				Content: r.URL.Query().Get("content"),
			}

			if err := msgPipeline.Execute(r.Context(), msg); err != nil {
				http.Error(w, "Failed to process message", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "Message sent to tenant: %s", tc.TenantID)
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
