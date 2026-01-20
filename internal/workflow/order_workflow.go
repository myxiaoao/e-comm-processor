package workflow

import (
	"log"
	"time"

	"go.temporal.io/sdk/workflow"

	"ecommerce-processor/internal/domain"
)

func OrderWorkflow(ctx workflow.Context, order domain.Order) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var payResult, invResult, refundResult, restockResult domain.ServiceResponse

	// 1. 支付处理
	err := workflow.ExecuteActivity(ctx, CallNatsService, domain.NatsSubjectPay, order).Get(ctx, &payResult)
	if err != nil || !payResult.Success {
		return "Payment failed", err
	}

	// 2. 库存预留
	err = workflow.ExecuteActivity(ctx, CallNatsService, domain.NatsSubjectInv, order).Get(ctx, &invResult)
	if err != nil || !invResult.Success {
		return "Inventory reservation failed", err
	}

	// 3. 检查是否需要取消（Saga 补偿逻辑）
	if order.Item == "Golang Book" {
		log.Printf("Cancellation request received for order %s", order.OrderID)

		err = workflow.ExecuteActivity(ctx, CallNatsService, domain.NatsSubjectRefund, order).Get(ctx, &refundResult)
		if err != nil {
			return "Refund failed (Critical Error)", err
		}

		err = workflow.ExecuteActivity(ctx, CallNatsService, domain.NatsSubjectRestock, order).Get(ctx, &restockResult)
		if err != nil {
			return "Restock failed (Critical Error)", err
		}

		return "Order Cancelled & Refunded (Saga Complete)", nil
	}

	return "Order Processed Successfully", nil
}
