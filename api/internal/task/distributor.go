package task

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskOrderConfirmation(
		ctx context.Context,
		payload OrderConfirmationPayload,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}

func (d *RedisTaskDistributor) DistributeTaskOrderConfirmation(
	ctx context.Context,
	payload OrderConfirmationPayload,
	opts ...asynq.Option,
) error {
	task, err := NewOrderConfirmationTask(payload.OrderID, payload.UserID, payload.TotalAmount)
	if err != nil {
		return err
	}

	_, err = d.client.EnqueueContext(ctx, task, opts...)
	return err
}
