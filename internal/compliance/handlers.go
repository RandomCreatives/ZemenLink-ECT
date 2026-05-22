package compliance

import (
	"encoding/json"
	"net/http"
	"zemenlink/internal/kernel"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Engine *WorkflowEngine
}

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/compliance", func(r chi.Router) {
		r.Use(h.RoleGuard("admin", "compliance_officer"))
		r.Post("/escrow-requests", h.InitiateEscrowRequest)
		r.Post("/escrow-requests/{id}/approve", h.ApproveEscrowRequest)
		r.Get("/audit-logs", h.GetAuditLogs)
	})
}

func (h *Handlers) RoleGuard(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value("user_role").(string)
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

func (h *Handlers) InitiateEscrowRequest(w http.ResponseWriter, r *http.Request) {
	tc, ok := kernel.GetTenantContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		TargetType string `json:"target_type"`
		TargetID   string `json:"target_id"`
		Reason     string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	req, err := h.Engine.InitiateRequest(r.Context(), tc.UserID, payload.TargetType, payload.TargetID, payload.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *Handlers) ApproveEscrowRequest(w http.ResponseWriter, r *http.Request) {
	tc, ok := kernel.GetTenantContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requestID := chi.URLParam(r, "id")
	if requestID == "" {
		http.Error(w, "Request ID required", http.StatusBadRequest)
		return
	}

	err := h.Engine.ApproveRequest(r.Context(), requestID, tc.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	tc, ok := kernel.GetTenantContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var logs []map[string]interface{}
	err := tc.DB.SelectContext(r.Context(), &logs, "SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
