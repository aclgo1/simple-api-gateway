package subscription

import (
	"context"
	"errors"
	"time"
)

type SubscriptionUseCase interface {
	CancelSubscription(ctx context.Context, params *ParamsCancelSubscriptionInput) (*ParamsCancelSubscriptionOutput, error)
}

type ParamsCancelSubscriptionInput struct {
	UserId string
}

func (p *ParamsCancelSubscriptionInput) Validate() error {
	if p.UserId == "" {
		return errors.New("user id empty")
	}
	return nil
}

type ParamsCancelSubscriptionOutput struct {
	SubscriptionId string    `json:"subscription_id"`
	UserId         string    `json:"user_id"`
	Status         string    `json:"status"`
	UpdatedAt      time.Time `json:"updated_at"`
}
