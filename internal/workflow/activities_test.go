package workflow

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-processor/internal/domain"
	"ecommerce-processor/internal/nats"
)

func startTestNatsServer() (*server.Server, error) {
	opts := &server.Options{
		Host: "127.0.0.1",
		Port: -1,
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	go ns.Start()
	if !ns.ReadyForConnections(5e9) {
		return nil, err
	}
	return ns, nil
}

func TestActivities_CallNatsService_Success(t *testing.T) {
	ns, err := startTestNatsServer()
	if err != nil {
		t.Skip("NATS server not available")
	}
	defer ns.Shutdown()

	natsClient, err := nats.NewClient(ns.ClientURL(), 2*time.Second)
	require.NoError(t, err)
	defer natsClient.Close()

	// 注册模拟服务
	_, err = natsClient.Conn().Subscribe(domain.NatsSubjectPay, func(m *natsgo.Msg) {
		resp := domain.ServiceResponse{Success: true, Message: "Payment Processed"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})
	require.NoError(t, err)

	activities := NewActivities(natsClient)
	order := domain.Order{OrderID: "test-001", Amount: 100.0, Item: "Test Item"}

	resp, err := activities.CallNatsService(context.Background(), domain.NatsSubjectPay, order)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Payment Processed", resp.Message)
}

func TestActivities_CallNatsService_Failure(t *testing.T) {
	ns, err := startTestNatsServer()
	if err != nil {
		t.Skip("NATS server not available")
	}
	defer ns.Shutdown()

	natsClient, err := nats.NewClient(ns.ClientURL(), 2*time.Second)
	require.NoError(t, err)
	defer natsClient.Close()

	// 注册返回失败的模拟服务
	_, err = natsClient.Conn().Subscribe(domain.NatsSubjectPay, func(m *natsgo.Msg) {
		resp := domain.ServiceResponse{Success: false, Message: "Insufficient Funds"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})
	require.NoError(t, err)

	activities := NewActivities(natsClient)
	order := domain.Order{OrderID: "test-002", Amount: 9999.0, Item: "Expensive Item"}

	resp, err := activities.CallNatsService(context.Background(), domain.NatsSubjectPay, order)

	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Insufficient Funds", err.Error())
}
