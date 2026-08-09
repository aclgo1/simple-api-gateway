package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/domain/models"
	protoUser "github.com/aclgo/simple-api-gateway/proto-service/user"
)


type subscriptionUC struct {
	subscriptionGRPC protoUser.SubscriptionServiceClient
}

func NewSubscriprionUseCase(subscriptionGRPC protoUser.SubscriptionServiceClient) *subscriptionUC{
	return &subscriptionUC{
		subscriptionGRPC: subscriptionGRPC,
	}
}

func(u *subscriptionUC)	ActivateSubscription(ctx context.Context, params *models.ParamsActivateSubscriptionInput)(*models.ParamsActivateSubscriptionOutput,error){
		pc := protoUser.CreateOrExtendSubscriptionRequest{
		UserId: params.AccountID,
		Plan: params.Plan,
		Days: int64(params.Days),
	}

	activate, err := u.subscriptionGRPC.CreateOrExtend(ctx, &pc)
	if err != nil {
		return nil, fmt.Errorf("u.subscriptionGRPC.CreateOrExtend: %w",err)
	}

	var expAt time.Time

	if activate.ExpiresAt != nil {
		expAt = activate.ExpiresAt.AsTime()
	}

	out := models.ParamsActivateSubscriptionOutput{
		SubscriptionID: activate.Id,
		ExpiresAt: expAt,
		Status: activate.Status,
	}

	return &out,nil
}