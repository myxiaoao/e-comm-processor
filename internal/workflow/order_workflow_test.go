package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"

	"ecommerce-processor/internal/domain"
)

func TestOrderWorkflow_Success(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterActivity(activityRef.CallNatsService)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectPay, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Payment OK"}, nil)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectInv, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Inventory OK"}, nil)

	order := domain.Order{
		OrderID: "001",
		Amount:  50.0,
		Item:    "Normal Item",
	}

	env.ExecuteWorkflow(OrderWorkflow, order)

	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())

	var result string
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Order Processed Successfully", result)
}

func TestOrderWorkflow_PaymentFailed(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterActivity(activityRef.CallNatsService)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectPay, mock.Anything).
		Return(domain.ServiceResponse{Success: false, Message: "Declined"}, nil)

	order := domain.Order{
		OrderID: "002",
		Amount:  100.0,
		Item:    "Test Item",
	}

	env.ExecuteWorkflow(OrderWorkflow, order)

	assert.True(t, env.IsWorkflowCompleted())

	var result string
	env.GetWorkflowResult(&result)
	assert.Equal(t, "Payment failed", result)
}

func TestOrderWorkflow_CancellationSaga(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterActivity(activityRef.CallNatsService)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectPay, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Payment OK"}, nil)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectInv, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Inventory OK"}, nil)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectRefund, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Refund OK"}, nil)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectRestock, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Restock OK"}, nil)

	order := domain.Order{
		OrderID: "003",
		Amount:  99.50,
		Item:    "Golang Book",
	}

	env.ExecuteWorkflow(OrderWorkflow, order)

	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())

	var result string
	assert.NoError(t, env.GetWorkflowResult(&result))
	assert.Equal(t, "Order Cancelled & Refunded (Saga Complete)", result)
}

func TestOrderWorkflow_InventoryFailed(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	env.RegisterActivity(activityRef.CallNatsService)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectPay, mock.Anything).
		Return(domain.ServiceResponse{Success: true, Message: "Payment OK"}, nil)
	env.OnActivity(activityRef.CallNatsService, mock.Anything, domain.NatsSubjectInv, mock.Anything).
		Return(domain.ServiceResponse{Success: false, Message: "Out of Stock"}, nil)

	order := domain.Order{
		OrderID: "004",
		Amount:  75.0,
		Item:    "Rare Item",
	}

	env.ExecuteWorkflow(OrderWorkflow, order)

	assert.True(t, env.IsWorkflowCompleted())

	var result string
	env.GetWorkflowResult(&result)
	assert.Equal(t, "Inventory reservation failed", result)
}
