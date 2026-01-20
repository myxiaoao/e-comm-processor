package nats

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ecommerce-processor/internal/domain"
)

func TestClient_Request_Success(t *testing.T) {
	// 启动内嵌 NATS 服务器
	ns, err := startTestServer()
	if err != nil {
		t.Skip("NATS server not available, skipping integration test")
	}
	defer ns.Shutdown()

	client, err := NewClient(ns.ClientURL(), 2*time.Second)
	require.NoError(t, err)
	defer client.Close()

	// 注册模拟响应
	_, err = client.Conn().Subscribe("test.subject", func(m *nats.Msg) {
		resp := domain.ServiceResponse{Success: true, Message: "OK"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})
	require.NoError(t, err)

	order := domain.Order{OrderID: "123", Amount: 50.0, Item: "Test Item"}
	resp, err := client.Request("test.subject", order)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "OK", resp.Message)
}

func TestClient_Request_Failure(t *testing.T) {
	ns, err := startTestServer()
	if err != nil {
		t.Skip("NATS server not available, skipping integration test")
	}
	defer ns.Shutdown()

	client, err := NewClient(ns.ClientURL(), 2*time.Second)
	require.NoError(t, err)
	defer client.Close()

	// 注册失败响应
	_, err = client.Conn().Subscribe("test.failure", func(m *nats.Msg) {
		resp := domain.ServiceResponse{Success: false, Message: "Payment Declined"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})
	require.NoError(t, err)

	order := domain.Order{OrderID: "456", Amount: 100.0, Item: "Expensive Item"}
	resp, err := client.Request("test.failure", order)

	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Payment Declined", resp.Message)
}

func TestClient_Request_Timeout(t *testing.T) {
	ns, err := startTestServer()
	if err != nil {
		t.Skip("NATS server not available, skipping integration test")
	}
	defer ns.Shutdown()

	client, err := NewClient(ns.ClientURL(), 100*time.Millisecond)
	require.NoError(t, err)
	defer client.Close()

	// 不注册订阅者，让请求超时
	order := domain.Order{OrderID: "789", Amount: 25.0, Item: "Timeout Item"}
	_, err = client.Request("nonexistent.subject", order)

	assert.Error(t, err)
}

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient("nats://invalid-host:9999", 100*time.Millisecond)
	assert.Error(t, err)
}
