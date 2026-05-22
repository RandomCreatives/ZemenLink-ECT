package messaging

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"zemenlink/internal/kernel"
)

type mockResolver struct {
	getDBFunc func(ctx context.Context, tenantID string) (interface{}, error)
}

func TestMessagingFlow(t *testing.T) {
	jwtSecret := []byte("test-secret")
	tenantID := "tenant-123"

	// 1. Generate Token
	token, err := kernel.GenerateTestToken(jwtSecret, tenantID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 2. Setup Pipeline
	pipeline := kernel.NewPipeline(func(ctx context.Context, msg *kernel.Message) error {
		msg.Metadata = map[string]interface{}{"processed": true}
		return nil
	})

	// 3. Setup Handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tc, ok := kernel.GetTenantContext(r.Context())
		if !ok {
			t.Error("Tenant context missing")
			return
		}
		if tc.TenantID != tenantID {
			t.Errorf("Expected tenant %s, got %s", tenantID, tc.TenantID)
		}

		msg := &kernel.Message{Content: "Hello"}
		if err := pipeline.Execute(r.Context(), msg); err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}

		if !msg.Metadata["processed"].(bool) {
			t.Error("Message was not processed by pipeline")
		}

		w.WriteHeader(http.StatusOK)
	})

	// 4. Run Request through Middleware
	// Note: We bypass the actual DB resolution for this unit test by injecting the context manually or mocking the resolver
	tc := &kernel.TenantContext{TenantID: tenantID}

	req := httptest.NewRequest("GET", "/messages", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a wrapper that injects the TC to avoid needing a real DB for this test
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(kernel.WithTenantContext(r.Context(), tc)))
		})
	}

	rr := httptest.NewRecorder()
	mw(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}
}
