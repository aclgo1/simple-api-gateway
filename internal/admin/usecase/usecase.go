package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/admin"
	"github.com/aclgo/simple-api-gateway/internal/auth"
	"github.com/aclgo/simple-api-gateway/internal/captcha"
	"github.com/aclgo/simple-api-gateway/pkg/logger"
	protoAdmin "github.com/aclgo/simple-api-gateway/proto-service/admin"
	protoUser "github.com/aclgo/simple-api-gateway/proto-service/user"

	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	protoMail "github.com/aclgo/simple-api-gateway/proto-service/mail"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type adminUC struct {
	clientAdmin   protoAdmin.AdminServiceClient
	clientUser    protoUser.UserServiceClient
	clientMail    protoMail.MailServiceClient
	clientBalance protoBalance.WalletServiceClient
	captchaRepo   captcha.Repository
	redisClient   *redis.Client
	logger        logger.Logger
}

func NewadminUC(clientAdmin protoAdmin.AdminServiceClient,
	clientUser protoUser.UserServiceClient,
	clientMail protoMail.MailServiceClient,
	clientBalance protoBalance.WalletServiceClient,
	captchaRepo captcha.Repository,
	redisClient *redis.Client,
	logger logger.Logger) admin.AdminUC {
	return &adminUC{
		clientAdmin:   clientAdmin,
		clientUser:    clientUser,
		clientMail:    clientMail,
		clientBalance: clientBalance,
		captchaRepo:   captchaRepo,
		redisClient:   redisClient,
		logger:        logger,
	}
}

func (u *adminUC) Create(ctx context.Context, params *admin.ParamsCreateAdmin) (*admin.Admin, error) {

	// if !u.captchaRepo.Verify(params.CaptchaId, params.CaptchaAwnser, true) {
	// 	return nil, admin.ErrFailedVerifyCaptcha{}
	// }

	in := protoAdmin.ParamsCreateAdmin{
		Name:     params.Name,
		Lastname: params.Lastname,
		Password: params.Password,
		Email:    params.Email,
		Role:     params.Role,
		Verified: params.Verified,
	}

	created, err := u.clientAdmin.Register(ctx, &in)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%+v\n", created)

	paramProtoBalance := protoBalance.ParamCreateWalletRequest{
		AccountID: created.UserId,
	}

	newWallet, err := u.clientBalance.Create(ctx, &paramProtoBalance)
	if err != nil {
		return nil, fmt.Errorf("u.clientBalance.Create: %w", err)
	}

	var newBalance float64
	if params.Balance > 0 {
		iw := protoBalance.ParamCreditWalletRequest{
			WalletID: newWallet.WalletID,
			Amount:   params.Balance,
		}

		c, err := u.clientBalance.Credit(ctx, &iw)
		if err != nil {
			u.logger.Error("failed credit user balance", err)
		}

		newBalance = c.Balance

	}

	if params.Verified == admin.DefaultVerfiedNo {
		err = u.redisClient.Get(ctx, created.Email).Err()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		confirmID := uuid.NewString()

		m := protoMail.MailRequest{
			From:        admin.DefaultFromSendMail,
			To:          params.Email,
			Subject:     admin.DefaultSubjectSendConfirm,
			Body:        fmt.Sprintf(admin.DefaulfBodySendConfirm, confirmID),
			Template:    admin.DefaulfTemplateSendConfirm,
			Servicename: admin.DefaultServiceSendName,
		}

		_, err = u.clientMail.SendService(ctx, &m)
		if err != nil {
			return nil, err
		}

		if err := u.redisClient.Set(ctx, params.Email, confirmID, time.Hour).Err(); err != nil {
			return nil, err
		}

		if err := u.redisClient.Set(ctx, confirmID, params.Email, time.Hour).Err(); err != nil {
			return nil, err
		}
	}

	return &admin.Admin{
		UserID:    created.UserId,
		Name:      created.Name,
		Lastname:  created.Lastname,
		Password:  created.Password,
		Email:     created.Email,
		Role:      created.Role,
		Verified:  created.Verified,
		Balance:   newBalance,
		CreatedAt: created.CreatedAt.AsTime(),
		UpdatedAt: created.UpdatedAt.AsTime(),
	}, nil
}

func (u *adminUC) Search(ctx context.Context, params *admin.ParamsSearch) ([]*admin.Admin, error) {

	in := protoAdmin.ParamsSearchRequest{
		Query:  params.Query,
		Role:   params.Role,
		Page:   int32(params.Page),
		Offset: int32(params.OffSet),
		Limit:  int32(params.Limit),
	}

	users, err := u.clientAdmin.Search(ctx, &in)
	if err != nil {
		return nil, err
	}

	items := make([]*admin.Admin, len(users.Users))

	var errs []error

	for i := 0; i < int(users.Total); i++ {
		paramWallet := protoBalance.ParamGetWalletByAccountRequest{
			AccountID: users.Users[i].UserId,
		}

		wallet, err := u.clientBalance.GetWalletByAccount(ctx, &paramWallet)
		if err != nil {
			errs = append(errs, err)
		}

		items[i] = &admin.Admin{
			UserID:    users.Users[i].UserId,
			Name:      users.Users[i].Name,
			Lastname:  users.Users[i].Lastname,
			Password:  users.Users[i].Password,
			Email:     users.Users[i].Email,
			Role:      users.Users[i].Role,
			Balance:   wallet.Balance,
			Verified:  users.Users[i].Verified,
			CreatedAt: users.Users[i].CreatedAt.AsTime(),
			UpdatedAt: users.Users[i].UpdatedAt.AsTime(),
		}
	}

	return items, errors.Join(errs...)
}

func (u *adminUC) Delete(ctx context.Context, params *admin.ParamsDeleteUser) (string, error) {

	pf := protoUser.FindByIdRequest{
		Id: params.UserId,
	}

	user, err := u.clientUser.FindById(ctx, &pf)

	if user.User.Role == "" {
		user.User.Role = string(auth.CLIENT)
	}

	if user.User.Role != string(auth.CLIENT) {
		return "", errors.New("user role unsuported delete")
	}

	i := protoAdmin.ParamsDeleteUserRequest{
		UserId: params.UserId,
	}

	o, err := u.clientAdmin.DeleteUser(ctx, &i)
	if err != nil {
		return "", err
	}

	return o.Msg, nil
}
