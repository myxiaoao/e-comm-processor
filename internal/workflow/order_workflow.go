package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"ecommerce-processor/internal/domain"
)

// activityRef is used to reference activity methods in workflow.ExecuteActivity calls.
// Temporal resolves the actual implementation from the registered Activities struct.
var activityRef *Activities

// OrderWorkflow processes an e-commerce order through payment, inventory, and optional cancellation steps.
func OrderWorkflow(ctx workflow.Context, order domain.Order) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var payResult, invResult domain.ServiceResponse

	// Step 1: Process payment
	err := workflow.ExecuteActivity(ctx, activityRef.CallNatsService, domain.NatsSubjectPay, order).Get(ctx, &payResult)
	if err != nil || !payResult.Success {
		return "Payment failed", err
	}

	// Step 2: Reserve inventory
	err = workflow.ExecuteActivity(ctx, activityRef.CallNatsService, domain.NatsSubjectInv, order).Get(ctx, &invResult)
	if err != nil || !invResult.Success {
		return "Inventory reservation failed", err
	}

	// Step 3: Check for cancellation (Saga compensation pattern)
	if order.Item == "Golang Book" {
		workflow.GetLogger(ctx).Info("Cancellation request received", "orderID", order.OrderID)
		return executeCancellation(ctx, order)
	}

	return "Order Processed Successfully", nil
}

// executeCancellation handles the Saga compensation: refund and restock.
func executeCancellation(ctx workflow.Context, order domain.Order) (string, error) {
	var refundResult, restockResult domain.ServiceResponse

	err := workflow.ExecuteActivity(ctx, activityRef.CallNatsService, domain.NatsSubjectRefund, order).Get(ctx, &refundResult)
	if err != nil {
		return "Refund failed (Critical Error)", err
	}

	err = workflow.ExecuteActivity(ctx, activityRef.CallNatsService, domain.NatsSubjectRestock, order).Get(ctx, &restockResult)
	if err != nil {
		return "Restock failed (Critical Error)", err
	}

	return "Order Cancelled & Refunded (Saga Complete)", nil
}
