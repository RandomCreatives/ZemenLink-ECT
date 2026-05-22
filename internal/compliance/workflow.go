package compliance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zemenlink/internal/kernel"
)

type EscrowRequest struct {
	ID                string    `db:"id"`
	RequesterID       string    `db:"requester_id"`
	TargetType        string    `db:"target_type"`
	TargetID          string    `db:"target_id"`
	Reason            string    `db:"reason"`
	Status            string    `db:"status"`
	ApprovalsRequired int       `db:"approvals_required"`
	CreatedAt         time.Time `db:"created_at"`
	ExpiresAt         time.Time `db:"expires_at"`
}

type WorkflowEngine struct{}

func (e *WorkflowEngine) InitiateRequest(ctx context.Context, requesterID, targetType, targetID, reason string) (*EscrowRequest, error) {
	tc, ok := kernel.GetTenantContext(ctx)
	if !ok {
		return nil, errors.New("tenant context missing")
	}

	req := &EscrowRequest{
		RequesterID:       requesterID,
		TargetType:        targetType,
		TargetID:          targetID,
		Reason:            reason,
		Status:            "pending",
		ApprovalsRequired: 2,
		CreatedAt:         time.Now(),
		ExpiresAt:         time.Now().Add(24 * time.Hour),
	}

	query := `
		INSERT INTO escrow_requests (requester_id, target_type, target_id, reason, status, approvals_required, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := tc.DB.QueryRowContext(ctx, query,
		req.RequesterID, req.TargetType, req.TargetID, req.Reason, req.Status, req.ApprovalsRequired, req.ExpiresAt,
	).Scan(&req.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to initiate escrow request: %w", err)
	}

	// Log the action
	e.logAction(ctx, requesterID, "ESCROW_REQUEST_INITIATED", "escrow_requests", req.ID)

	return req, nil
}

func (e *WorkflowEngine) ApproveRequest(ctx context.Context, requestID, approverID string) error {
	tc, ok := kernel.GetTenantContext(ctx)
	if !ok {
		return errors.New("tenant context missing")
	}

	tx, err := tc.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Record the approval
	_, err = tx.ExecContext(ctx, "INSERT INTO escrow_approvals (request_id, approver_id) VALUES ($1, $2)", requestID, approverID)
	if err != nil {
		return fmt.Errorf("failed to record approval: %w", err)
	}

	// 2. Check if threshold is met
	var approvalCount int
	err = tx.GetContext(ctx, &approvalCount, "SELECT count(*) FROM escrow_approvals WHERE request_id = $1", requestID)
	if err != nil {
		return err
	}

	var req EscrowRequest
	err = tx.GetContext(ctx, &req, "SELECT * FROM escrow_requests WHERE id = $1", requestID)
	if err != nil {
		return err
	}

	if approvalCount >= req.ApprovalsRequired {
		_, err = tx.ExecContext(ctx, "UPDATE escrow_requests SET status = 'approved' WHERE id = $1", requestID)
		if err != nil {
			return err
		}
		// Trigger ReleaseKey event logic here
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	e.logAction(ctx, approverID, "ESCROW_REQUEST_APPROVED", "escrow_requests", requestID)
	return nil
}

func (e *WorkflowEngine) logAction(ctx context.Context, userID, action, resType, resID string) {
	tc, ok := kernel.GetTenantContext(ctx)
	if !ok {
		return
	}
	_, _ = tc.DB.ExecContext(ctx,
		"INSERT INTO audit_logs (user_id, action, resource_type, resource_id) VALUES ($1, $2, $3, $4)",
		userID, action, resType, resID)
}
