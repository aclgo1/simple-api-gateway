package models

import "time"

type ParamsActivateSubscriptionInput struct {
	AccountID string
	Plan      string
	Days      int
}

type ParamsActivateSubscriptionOutput struct {
	SubscriptionID string
	ExpiresAt      time.Time
	Status         string
}
