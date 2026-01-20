package workflow

import (
	"context"

	"ecommerce-processor/internal/domain"
	"ecommerce-processor/internal/nats"
)

// Activities holds dependencies for workflow activities.
type Activities struct {
	natsClient *nats.Client
}

// NewActivities creates a new Activities instance with the given NATS client.
func NewActivities(natsClient *nats.Client) *Activities {
	return &Activities{natsClient: natsClient}
}

// CallNatsService sends a request to a NATS service and returns the response.
// This method is registered as a Temporal activity.
func (a *Activities) CallNatsService(ctx context.Context, subject string, order domain.Order) (domain.ServiceResponse, error) {
	return a.natsClient.Request(subject, order)
}
