package domain

const (
	TaskQueue          = "ORDER_QUEUE"
	NatsSubjectInv     = "service.inventory"
	NatsSubjectPay     = "service.payment"
	NatsSubjectRefund  = "service.refund"
	NatsSubjectRestock = "service.restock"
)

type Order struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
}

type ServiceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
