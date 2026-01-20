package workflow

import (
	"context"

	"ecommerce-processor/internal/domain"
	"ecommerce-processor/internal/nats"
)

type Activities struct {
	natsClient *nats.Client
}

func NewActivities(natsClient *nats.Client) *Activities {
	return &Activities{natsClient: natsClient}
}

func (a *Activities) CallNatsService(ctx context.Context, subject string, order domain.Order) (domain.ServiceResponse, error) {
	return a.natsClient.Request(subject, order)
}

// CallNatsService 独立函数，用于工作流中调用
func CallNatsService(ctx context.Context, subject string, order domain.Order) (domain.ServiceResponse, error) {
	// 这个函数由 Worker 注册时的 Activities 实例提供实现
	// 在工作流中只作为 Activity 类型引用
	return domain.ServiceResponse{}, nil
}
