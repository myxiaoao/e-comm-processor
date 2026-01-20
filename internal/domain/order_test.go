package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_JSONSerialization(t *testing.T) {
	order := Order{
		OrderID: "123",
		Amount:  99.50,
		Item:    "Golang Book",
	}

	data, err := json.Marshal(order)
	assert.NoError(t, err)

	var decoded Order
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, order, decoded)
}

func TestServiceResponse_JSONSerialization(t *testing.T) {
	resp := ServiceResponse{
		Success: true,
		Message: "Transaction Approved",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded ServiceResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, resp, decoded)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "ORDER_QUEUE", TaskQueue)
	assert.Equal(t, "service.inventory", NatsSubjectInv)
	assert.Equal(t, "service.payment", NatsSubjectPay)
	assert.Equal(t, "service.refund", NatsSubjectRefund)
	assert.Equal(t, "service.restock", NatsSubjectRestock)
}
