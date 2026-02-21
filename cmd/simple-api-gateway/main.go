package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aclgo/simple-api-gateway/config"
	"github.com/aclgo/simple-api-gateway/frontend/load"

	"github.com/aclgo/simple-api-gateway/internal/admin"
	captchaRepo "github.com/aclgo/simple-api-gateway/internal/captcha/repository"
	captchaUC "github.com/aclgo/simple-api-gateway/internal/captcha/usecase"
	svcAdmin "github.com/aclgo/simple-api-gateway/internal/delivery/http/service/admin"
	svcCaptcha "github.com/aclgo/simple-api-gateway/internal/delivery/http/service/captcha"
	svcOrders "github.com/aclgo/simple-api-gateway/internal/delivery/http/service/orders"
	svcProduct "github.com/aclgo/simple-api-gateway/internal/delivery/http/service/product"
	svcUser "github.com/aclgo/simple-api-gateway/internal/delivery/http/service/user"
	svcPix "github.com/aclgo/simple-api-gateway/internal/delivery/http/service/wallet/pix"
	svcEx "github.com/aclgo/simple-api-gateway/internal/delivery/websocket/service/ex"
	"github.com/aclgo/simple-api-gateway/internal/user"
	"github.com/aclgo/simple-api-gateway/internal/wallet"
	migration "github.com/aclgo/simple-api-gateway/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	adminUC "github.com/aclgo/simple-api-gateway/internal/admin/usecase"
	authUC "github.com/aclgo/simple-api-gateway/internal/auth/usecase"
	ordersUC "github.com/aclgo/simple-api-gateway/internal/orders/usecase"
	productUC "github.com/aclgo/simple-api-gateway/internal/product/usecase"
	userUC "github.com/aclgo/simple-api-gateway/internal/user/usecase"
	cardUC "github.com/aclgo/simple-api-gateway/internal/wallet/card/usecase"
	pixUC "github.com/aclgo/simple-api-gateway/internal/wallet/pix/usecase"
	walletUC "github.com/aclgo/simple-api-gateway/internal/wallet/usecase"

	redis "github.com/aclgo/simple-api-gateway/pkg/rredis"

	"github.com/aclgo/simple-api-gateway/pkg/logger"
	protoAdmin "github.com/aclgo/simple-api-gateway/proto-service/admin"
	protoBalance "github.com/aclgo/simple-api-gateway/proto-service/balance"
	protoMail "github.com/aclgo/simple-api-gateway/proto-service/mail"
	protoOrders "github.com/aclgo/simple-api-gateway/proto-service/orders"
	protoProduct "github.com/aclgo/simple-api-gateway/proto-service/product"
	protoUser "github.com/aclgo/simple-api-gateway/proto-service/user"
)

var (
	AddrServiceAdmin    = "grpc-admin:50051"
	OptionsServiceAdmin = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	AddrServiceUser    = "grpc-user:50052"
	OptionsServiceUser = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	AddrServiceMail    = "grpc-mail:50053"
	OptionsServiceMail = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	AddrServiceProduct    = "grpc-product:50054"
	OptionsServiceProduct = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	AddrServiceOrders    = "grpc-orders:50055"
	OptionsServiceOrders = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	AddrServiceBalance    = "grpc-balance:50056"
	OptionsServiceBalance = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
)

func main() {

	cfg := config.Load(".")

	db, err := sqlx.Open(cfg.DbDriver, cfg.DbUrl)
	if err != nil {
		log.Fatalf("sqlx.Open: %v", err)
	}

	migration.NewMigration(db, nil)
	if err := migration.Run(); err != nil {
		log.Fatalf("migration.Run: %v", err)
	}

	logger, err := logger.NewapiLogger(cfg)
	if err != nil {
		log.Fatalf("logger.NewapiLogger: %s\n", err)
	}

	logger.Info("logger initialized")

	//	CONNECTING IN MICROSERVICES
	connUser, err := grpc.NewClient(AddrServiceUser, OptionsServiceUser...)
	if err != nil {
		logger.Errorf("grpc.Dial: connection in user service: %v", err)
	}

	connAdmin, err := grpc.NewClient(AddrServiceAdmin, OptionsServiceAdmin...)
	if err != nil {
		logger.Errorf("grpc.Dial: connection in admin service: %v", err)
	}

	connMail, err := grpc.NewClient(AddrServiceMail, OptionsServiceMail...)
	if err != nil {
		logger.Errorf("grpc.Dial: connection in mail service: %v", err)
	}

	connProduct, err := grpc.NewClient(AddrServiceProduct, OptionsServiceProduct...)
	if err != nil {
		logger.Errorf("grpc.Dial: connection in product service: %v", err)
	}

	connOrders, err := grpc.NewClient(AddrServiceOrders, OptionsServiceOrders...)
	if err != nil {
		logger.Errorf("grpc.Dial: connection in product service: %v", err)
	}

	connBalance, err := grpc.NewClient(AddrServiceBalance, OptionsServiceBalance...)
	if err != nil {
		logger.Errorf("grpc.Dial: connection in product service: %v", err)
	}

	redisClient := redis.NewRedisClient(cfg)

	mu := sync.Mutex{}
	user.SetConfigUserPackage(cfg.BaseApiUrl, cfg.DefaultEmailSendEmail, cfg.DefaultTimeSendEmail, cfg.DefaultServiceNameSendEmail)
	admin.SetConfigUserPackage(cfg.BaseApiUrl, cfg.DefaultEmailSendEmail, cfg.DefaultTimeSendEmail, cfg.DefaultServiceNameSendEmail)

	////////////////////////////////

	clientUserService := protoUser.NewUserServiceClient(connUser)
	adminUserService := protoAdmin.NewAdminServiceClient(connAdmin)
	mailUserService := protoMail.NewMailServiceClient(connMail)
	productUserService := protoProduct.NewProductServiceClient(connProduct)
	ordersUserService := protoOrders.NewServiceOrderClient(connOrders)
	balanceUserService := protoBalance.NewWalletServiceClient(connBalance)

	cptRepo := captchaRepo.NewRepository()
	cptUC := captchaUC.NewCaptchaUC(cptRepo)
	cptSvc := svcCaptcha.NewCaptchaService(cptUC)

	user := userUC.NewuserUC(clientUserService, mailUserService, balanceUserService, cptRepo, redisClient, logger)
	admin := adminUC.NewadminUC(adminUserService, clientUserService, mailUserService, balanceUserService, cptRepo, redisClient, logger)
	product := productUC.NewProductUC(logger, productUserService)
	orders := ordersUC.NeworderUC(ordersUserService, productUserService, balanceUserService, &mu, logger)

	w := walletUC.NewwalletUC(balanceUserService, logger)

	pixProcessor := pixUC.NewpaymentProcessorPix()
	cardProcessor := cardUC.NewpaymentProcessorCard()

	w.RegisterProvider(wallet.PaymentMethodPix, pixProcessor)
	w.RegisterProvider(wallet.PaymentMethodCard, cardProcessor)

	userHandler := svcUser.NewuserService(user, logger, cfg.BaseApiUrl)
	adminHandler := svcAdmin.NewadminService(admin, logger)
	productHandler := svcProduct.NewProductService(product, logger)
	ordersHandler := svcOrders.NewOrdersService(orders, logger)
	walletPixHandler := svcPix.NewwalletServicePix(pixProcessor, w)
	exHandler := svcEx.NewExService()

	authUC := authUC.NewAuthUC(clientUserService)

	ctx := context.Background()

	mux := http.NewServeMux()

	//MICROSERVICE GRPC USER
	mux.HandleFunc("POST /api/login", userHandler.Login(ctx))
	mux.HandleFunc("GET /api/logout", authUC.ValidateTwoToken(userHandler.Logout(ctx)))
	mux.HandleFunc("GET /api/refresh", authUC.ValidateTwoToken(userHandler.RefreshTokens(ctx)))
	mux.HandleFunc("GET /api/valid_token", authUC.ValidateToken(userHandler.ValidToken(ctx)))
	mux.HandleFunc("POST /api/user/register", userHandler.Register(ctx))
	mux.HandleFunc("GET /api/user/find", authUC.ValidateToken(userHandler.Find(ctx)))
	mux.HandleFunc("PUT /api/user/update", authUC.ValidateUpdate(userHandler.Update(ctx)))

	//MICROSERVCE GRPC MAIL
	mux.HandleFunc("GET /api/user/confirm/{confirm_code}", userHandler.UserConfirm(ctx))
	mux.HandleFunc("GET /api/user/resetpass/{email}/{captcha_id}/{captcha_awnser}", userHandler.UserResetPass(ctx))
	mux.HandleFunc("POST /api/user/newpass/{code}", userHandler.UserNewPass(ctx))

	//MICROSERVICE GRPC ADMINs
	mux.HandleFunc("/api/admin/register", authUC.ValidateCreateAdmin(adminHandler.Create(ctx)))
	mux.HandleFunc("/api/admin/search", authUC.ValidateIsAdmin(adminHandler.Search(ctx)))
	mux.HandleFunc("DELETE /api/admin/delete/{user_id}", authUC.ValidateIsAdmin(adminHandler.Delete(ctx)))

	//MICROSERVICE GRPC PRODUCTS
	mux.HandleFunc("POST /api/product/create", authUC.ValidateIsAdmin(productHandler.Create(ctx)))
	mux.HandleFunc("GET /api/product/find/{product_id}", authUC.ValidateToken(productHandler.Find(ctx)))
	mux.HandleFunc("GET /api/product/findall", productHandler.FindAll(ctx))
	mux.HandleFunc("PUT /api/product/update", authUC.ValidateIsAdmin(productHandler.Update(ctx)))
	mux.HandleFunc("DELETE /api/product/delete/{product_id}", authUC.ValidateIsAdmin(productHandler.Delete(ctx)))

	//MICROSERVICE GRPC ORDERS
	mux.HandleFunc("POST /api/orders", authUC.ValidateToken(ordersHandler.Create(ctx)))
	mux.HandleFunc("GET /api/orders/find/{order_id}", authUC.ValidateIsAdmin(ordersHandler.FindById(ctx)))
	mux.HandleFunc("GET /api/orders/find/account", authUC.ValidateToken(ordersHandler.FindByAccount(ctx)))
	mux.HandleFunc("GET /api/orders/find/product/{product_id}", authUC.ValidateIsAdmin(ordersHandler.FindByProduct(ctx)))

	mux.HandleFunc("POST /api/payments/pix", authUC.ValidateToken(walletPixHandler.CreatePix(ctx)))
	mux.HandleFunc("POST /api/webhook/pix", authUC.ValidateWebHookPix(ctx, walletPixHandler.WebHookPix()))

	mux.HandleFunc("GET /api/captcha", cptSvc.GenCaptcha(ctx))

	mux.HandleFunc("GET /api/ws/card/{token}", authUC.ValidateTokenWs(exHandler.ExWs(ctx)))
	mux.HandleFunc("GET /api/ws/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if e := json.NewEncoder(w).Encode([]map[string]any{{"path": "/api/ws/card/"}}); e != nil {
			log.Printf("json.NewEncoder: %v\n", e)
		}
	})

	//FRONTEND SETUP

	load := load.NewLoad("html", "css", "./")

	pages, err := load.Start()
	if err != nil {
		log.Fatalf("load.Start: %v", err)
	}

	// mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(load.PathCss))))

	mux.HandleFunc("/login", pages.Login)
	mux.HandleFunc("/home", pages.Home)
	mux.HandleFunc("/unauthorized", pages.Unauthorized)
	mux.HandleFunc("/confirm_signup", pages.ConfirmSignup)
	mux.HandleFunc("/resetpass", pages.ResetPass)
	mux.HandleFunc("/newpass", pages.NewPass)
	mux.HandleFunc("/products", pages.Products)
	mux.HandleFunc("/admin", pages.Admin)
	mux.HandleFunc("/ws", pages.Ws)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
	})

	hlogger := logger.BasicHttpLogger(mux)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.ApiPort),
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		ErrorLog:       log.Default(),
		Handler:        cors.Handler(hlogger),
		MaxHeaderBytes: 8192,
	}

	logger.Infof("server running port %d", cfg.ApiPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("mux.ListenAndServe:%v", err)
	}
}
