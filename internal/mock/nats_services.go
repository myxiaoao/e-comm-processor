package mock

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"

	"ecommerce-processor/internal/domain"
)

type mockService struct {
	subject     string
	logMessage  func(order domain.Order) string
	respMessage string
}

func StartNatsServices(nc *nats.Conn) {
	services := []mockService{
		{
			subject:     domain.NatsSubjectInv,
			logMessage:  func(o domain.Order) string { return fmt.Sprintf("[Mock Inventory] Checking inventory for: %s", o.Item) },
			respMessage: "Item Reserved",
		},
		{
			subject:     domain.NatsSubjectPay,
			logMessage:  func(o domain.Order) string { return fmt.Sprintf("[Mock Payment] Charging credit card: $%.2f", o.Amount) },
			respMessage: "Transaction Approved",
		},
		{
			subject:     domain.NatsSubjectRefund,
			logMessage:  func(o domain.Order) string { return fmt.Sprintf("[Mock Refund] Issuing refund for order: %s ($%.2f)", o.OrderID, o.Amount) },
			respMessage: "Refund Processed",
		},
		{
			subject:     domain.NatsSubjectRestock,
			logMessage:  func(o domain.Order) string { return fmt.Sprintf("[Mock Restock] Restocking item: %s", o.Item) },
			respMessage: "Item Restocked",
		},
	}

	for _, svc := range services {
		registerMockHandler(nc, svc)
	}

	log.Println("Mock NATS services started (inventory, payment, refund, restock)")
}

func registerMockHandler(nc *nats.Conn, svc mockService) {
	nc.Subscribe(svc.subject, func(m *nats.Msg) {
		var order domain.Order
		json.Unmarshal(m.Data, &order)

		log.Print(svc.logMessage(order))

		resp := domain.ServiceResponse{Success: true, Message: svc.respMessage}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})
}
