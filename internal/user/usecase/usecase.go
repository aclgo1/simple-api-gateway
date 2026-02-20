package usecase

import (
	"context"
	"fmt"

	"github.com/aclgo/simple-api-gateway/internal/captcha"
	"github.com/aclgo/simple-api-gateway/internal/user"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	mail "github.com/aclgo/simple-api-gateway/proto-service/mail"
	protoUser "github.com/aclgo/simple-api-gateway/proto-service/user"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type userUc struct {
	clientUserGRPC    protoUser.UserServiceClient
	clientMailGRPC    mail.MailServiceClient
	clientBalanceGRPC protoBalance.WalletServiceClient
	captchaRepo       captcha.Repository
	redisClient       *redis.Client
	baseApiUrl        string
	logger            logger.Logger
}

func NewuserUC(clientUser protoUser.UserServiceClient,
	clientMail mail.MailServiceClient,
	clientBalance protoBalance.WalletServiceClient,
	captchaRepo captcha.Repository,
	redisClient *redis.Client,
	logger logger.Logger,
) user.UserUC {

	return &userUc{
		clientUserGRPC:    clientUser,
		clientMailGRPC:    clientMail,
		clientBalanceGRPC: clientBalance,
		captchaRepo:       captchaRepo,
		redisClient:       redisClient,
		logger:            logger,
	}
}

func (u *userUc) Register(ctx context.Context, params *user.ParamsUserRegister) (*user.User, error) {
	if !u.captchaRepo.Verify(params.CaptchaId, params.CaptchaAwnser, true) {
		return nil, user.ErrFailedVerifyCaptcha{}
	}

	created, err := u.clientUserGRPC.Register(ctx, &protoUser.CreateUserRequest{
		Name:     params.Name,
		LastName: params.Lastname,
		Password: params.Password,
		Email:    params.Email,
	})

	if err != nil {
		return nil, err
	}

	paramCreateBalance := protoBalance.ParamCreateWalletRequest{
		AccountID: created.User.Id,
	}

	wallet, err := u.clientBalanceGRPC.Create(ctx, &paramCreateBalance)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalanceGRPC.Create: %w", err)
	}

	return &user.User{
		UserID:    created.User.Id,
		Name:      created.User.Name,
		Lastname:  created.User.LastName,
		Password:  created.User.Password,
		Email:     created.User.Email,
		Role:      created.User.Role,
		Verified:  created.User.Verified,
		Balance:   wallet.Balance,
		CreatedAt: created.User.CreatedAt.AsTime(),
		UpdatedAt: created.User.UpdatedAt.AsTime(),
	}, nil
}

func (u *userUc) Login(ctx context.Context, params *user.ParamsUserLoginRequest) (*user.ParamsUserLoginResponse, error) {
	if !u.captchaRepo.Verify(params.CaptchaId, params.CaptchaAwnser, true) {
		return nil, user.ErrFailedVerifyCaptcha{}
	}

	resp, err := u.clientUserGRPC.Login(ctx, &protoUser.UserLoginRequest{
		Email:    params.Email,
		Password: params.Password,
	})

	if err != nil && err != (user.ErrUserNotVerified{}) {
		return nil, err
	}

	if err == (user.ErrUserNotVerified{}) {
		return nil, user.ErrUserNotVerified{}
	}

	return &user.ParamsUserLoginResponse{
		AccessToken:  resp.Tokens.AccessToken,
		RefreshToken: resp.Tokens.RefreshToken,
	}, nil
}

func (u *userUc) Logout(ctx context.Context, params *user.ParamsUserLogout) error {
	_, err := u.clientUserGRPC.Logout(ctx, &protoUser.UserLogoutRequest{
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
	})

	if err != nil {
		return err
	}

	return nil
}
func (u *userUc) FindById(ctx context.Context, params *user.ParamsUserFindById) (*user.User, error) {
	resp, err := u.clientUserGRPC.FindById(ctx, &protoUser.FindByIdRequest{Id: params.UserID})
	if err != nil {
		return nil, err
	}

	paramFindWallet := protoBalance.ParamGetWalletByAccountRequest{
		AccountID: params.UserID,
	}

	wallet, err := u.clientBalanceGRPC.GetWalletByAccount(ctx, &paramFindWallet)
	if err != nil {
		u.logger.Error("u.clientBalanceGRPC.GetWalletByAccount", err)

		paramProtoNewWallet := protoBalance.ParamCreateWalletRequest{
			AccountID: params.UserID,
		}

		newWallet, err := u.clientBalanceGRPC.Create(ctx, &paramProtoNewWallet)
		if err != nil {
			u.logger.Error("u.clientBalanceGRPC.Create: ", err)
		}

		return &user.User{
			UserID:    resp.User.Id,
			Name:      resp.User.Name,
			Lastname:  resp.User.LastName,
			Password:  resp.User.Password,
			Email:     resp.User.Email,
			Role:      resp.User.Role,
			Verified:  resp.User.Verified,
			Balance:   newWallet.Balance,
			CreatedAt: resp.User.CreatedAt.AsTime(),
			UpdatedAt: resp.User.UpdatedAt.AsTime(),
		}, nil
	}

	return &user.User{
		UserID:    resp.User.Id,
		Name:      resp.User.Name,
		Lastname:  resp.User.LastName,
		Password:  resp.User.Password,
		Email:     resp.User.Email,
		Role:      resp.User.Role,
		Verified:  resp.User.Verified,
		Balance:   wallet.Balance,
		CreatedAt: resp.User.CreatedAt.AsTime(),
		UpdatedAt: resp.User.UpdatedAt.AsTime(),
	}, nil
}
func (u *userUc) FindByEmail(ctx context.Context, params *user.ParamsUserFindByEmail) (*user.User, error) {
	resp, err := u.clientUserGRPC.FindByEmail(ctx, &protoUser.FindByEmailRequest{Email: params.Email})
	if err != nil {
		return nil, err
	}

	paramProtoFind := protoBalance.ParamGetWalletByAccountRequest{
		AccountID: resp.User.Id,
	}

	wallet, err := u.clientBalanceGRPC.GetWalletByAccount(ctx, &paramProtoFind)
	if err != nil {
		u.logger.Error("u.clientBalanceGRPC.GetWalletByAccount", err)

		paramProtoNewWallet := protoBalance.ParamCreateWalletRequest{
			AccountID: resp.User.Id,
		}

		newWallet, err := u.clientBalanceGRPC.Create(ctx, &paramProtoNewWallet)
		if err != nil {
			u.logger.Error("u.clientBalanceGRPC.Create: ", err)
		}

		return &user.User{
			UserID:    resp.User.Id,
			Name:      resp.User.Name,
			Lastname:  resp.User.LastName,
			Password:  resp.User.Password,
			Email:     resp.User.Email,
			Role:      resp.User.Role,
			Verified:  resp.User.Verified,
			Balance:   newWallet.Balance,
			CreatedAt: resp.User.CreatedAt.AsTime(),
			UpdatedAt: resp.User.UpdatedAt.AsTime(),
		}, nil
	}

	return &user.User{
		UserID:    resp.User.Id,
		Name:      resp.User.Name,
		Lastname:  resp.User.LastName,
		Password:  resp.User.Password,
		Email:     resp.User.Email,
		Role:      resp.User.Role,
		Verified:  resp.User.Verified,
		Balance:   wallet.Balance,
		CreatedAt: resp.User.CreatedAt.AsTime(),
		UpdatedAt: resp.User.UpdatedAt.AsTime(),
	}, nil
}
func (u *userUc) Update(ctx context.Context, params *user.ParamsUserUpdate) (*user.User, error) {
	updated, err := u.clientUserGRPC.Update(ctx,
		&protoUser.UpdateRequest{
			Id:       params.UserID,
			Name:     params.Name,
			Lastname: params.Lastname,
			Password: params.Password,
			Email:    params.Email,
		},
	)
	if err != nil {
		return nil, err
	}

	paramProtoFind := protoBalance.ParamGetWalletByAccountRequest{
		AccountID: params.UserID,
	}

	wallet, err := u.clientBalanceGRPC.GetWalletByAccount(ctx, &paramProtoFind)
	if err != nil {
		return nil, err
	}

	newBalance := 0.0
	if params.Balance > 0 {
		paramProtoUpdateBalance := protoBalance.ParamCreditWalletRequest{
			WalletID: wallet.WalletID,
			Amount:   params.Balance,
		}

		walletUpdated, err := u.clientBalanceGRPC.Credit(ctx, &paramProtoUpdateBalance)
		if err != nil {
			return nil, err
		}

		newBalance = walletUpdated.Balance
	}

	return &user.User{
		UserID:    updated.User.Id,
		Name:      updated.User.Name,
		Lastname:  updated.User.LastName,
		Password:  updated.User.Password,
		Email:     updated.User.Email,
		Role:      updated.User.Role,
		Balance:   newBalance,
		Verified:  updated.User.Verified,
		CreatedAt: updated.User.CreatedAt.AsTime(),
		UpdatedAt: updated.User.UpdatedAt.AsTime(),
	}, nil
}

func (u *userUc) Delete(ctx context.Context, params *user.ParamsUserDelete) error {
	_, err := u.clientUserGRPC.Delete(ctx, &protoUser.DeleteRequest{Id: params.UserID})
	return err
}

func (u *userUc) RefreshTokens(ctx context.Context, params *user.ParamsRefreshTokens) (*user.RefreshTokens, error) {
	tokens, err := u.clientUserGRPC.RefreshTokens(ctx, &protoUser.RefreshTokensRequest{
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
	})

	if err != nil {
		return nil, err
	}

	return &user.RefreshTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (u *userUc) SendConfirm(ctx context.Context, params *user.ParamsConfirm) error {

	err := u.redisClient.Get(ctx, params.To).Err()
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {

		confirmID := uuid.NewString()

		req := &mail.MailRequest{
			From:        user.DefaultFromSendMail,
			To:          params.To,
			Subject:     user.DefaultSubjectSendConfirm,
			Body:        fmt.Sprintf(user.DefaulfBodySendConfirm, confirmID),
			Template:    user.DefaulfTemplateSendConfirm,
			Servicename: user.DefaultServiceName,
		}

		_, err = u.clientMailGRPC.SendService(ctx, req)
		if err != nil {
			return err
		}

		if err := u.redisClient.Set(ctx, params.To, confirmID, user.DefaultTimeSendEmails).Err(); err != nil {
			return err
		}

		if err := u.redisClient.Set(ctx, confirmID, params.To, user.DefaultTimeSendEmails).Err(); err != nil {
			return err
		}

		return nil
	}

	return user.ErrEmailSentCheckInbox{}
}

func (u *userUc) SendConfirmOK(ctx context.Context, params *user.ParamsConfirmOK) error {

	userEmail, err := u.redisClient.Get(ctx, params.ConfirmCode).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		return user.ErrInvalidCode{}
	}

	foundUser, err := u.clientUserGRPC.FindByEmail(ctx, &protoUser.FindByEmailRequest{Email: userEmail})

	if err != nil {
		return err
	}

	in := protoUser.UpdateRequest{
		Id:       foundUser.User.Id,
		Verified: "yes",
	}

	_, err = u.clientUserGRPC.Update(ctx, &in)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUc) ResetPass(ctx context.Context, params *user.ParamsResetPass) error {
	if !u.captchaRepo.Verify(params.CaptchaId, params.CaptchaAwnser, true) {
		return user.ErrFailedVerifyCaptcha{}
	}

	userData, err := u.clientUserGRPC.FindByEmail(ctx, &protoUser.FindByEmailRequest{Email: params.Email})
	if err != nil {
		return fmt.Errorf("u.clientUserGRPC.FindByEmail: %w", err)
	}

	if userData.User.Verified == user.DefaultVerifiedNo {
		return user.ErrUserNotVerified{}
	}

	errRedis := u.redisClient.Get(ctx, userData.User.Id).Err()
	if errRedis != redis.Nil && errRedis != nil {
		return errRedis
	}

	resetCode := uuid.NewString()

	if errRedis == redis.Nil {
		in := mail.MailRequest{
			From:        user.DefaultFromSendMail,
			To:          params.Email,
			Subject:     user.DefaultSubjectResetPass,
			Body:        fmt.Sprintf(user.DefaultBodyResetPass, resetCode),
			Template:    user.DefaultTemplateResetPass,
			Servicename: user.DefaultServiceName,
		}

		_, err = u.clientMailGRPC.SendService(ctx, &in)
		if err != nil {
			return err
		}

		if err := u.redisClient.Set(ctx, resetCode, userData.User.Id, user.DefaultTimeSendEmails).Err(); err != nil {
			return err
		}

		if err := u.redisClient.Set(ctx, userData.User.Id, resetCode, user.DefaultTimeSendEmails).Err(); err != nil {
			return err
		}

		return nil
	}

	return user.ErrEmailSentCheckInbox{}

}

func (u *userUc) NewPass(ctx context.Context, params *user.ParamsNewPass) error {
	if !u.captchaRepo.Verify(params.CaptchaId, params.CaptchaAwnser, true) {
		return user.ErrFailedVerifyCaptcha{}
	}

	idUser, err := u.redisClient.Get(ctx, params.NewPassCode).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		return user.ErrInvalidCode{}
	}

	_, err = u.clientUserGRPC.FindById(ctx, &protoUser.FindByIdRequest{Id: idUser})
	if err != nil {
		return err
	}

	_, err = u.clientUserGRPC.Update(ctx, &protoUser.UpdateRequest{
		Id:       idUser,
		Password: params.NewPass,
	})

	if err != nil {
		return err
	}

	errDel := u.redisClient.Del(ctx, params.NewPassCode).Err()
	if errDel != nil {
		return errDel
	}

	return nil
}
