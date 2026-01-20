package mock

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"ecommerce-processor/internal/domain"
)

func StartNatsServices(nc *nats.Conn) {
	nc.Subscribe(domain.NatsSubjectInv, func(m *nats.Msg) {
		var order domain.Order
		json.Unmarshal(m.Data, &order)

		log.Printf("[Mock Inventory] Checking inventory for: %s", order.Item)

		resp := domain.ServiceResponse{Success: true, Message: "Item Reserved"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})

	nc.Subscribe(domain.NatsSubjectPay, func(m *nats.Msg) {
		var order domain.Order
		json.Unmarshal(m.Data, &order)

		log.Printf("[Mock Payment] Charging credit card: $%.2f", order.Amount)

		resp := domain.ServiceResponse{Success: true, Message: "Transaction Approved"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})

	nc.Subscribe(domain.NatsSubjectRefund, func(m *nats.Msg) {
		var order domain.Order
		json.Unmarshal(m.Data, &order)

		log.Printf("[Mock Refund] Issuing refund for order: %s ($%.2f)", order.OrderID, order.Amount)

		resp := domain.ServiceResponse{Success: true, Message: "Refund Processed"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})

	nc.Subscribe(domain.NatsSubjectRestock, func(m *nats.Msg) {
		var order domain.Order
		json.Unmarshal(m.Data, &order)

		log.Printf("[Mock Restock] Restocking item: %s", order.Item)

		resp := domain.ServiceResponse{Success: true, Message: "Item Restocked"}
		data, _ := json.Marshal(resp)
		m.Respond(data)
	})

	log.Println("Mock NATS services started (inventory, payment, refund, restock)")
}
