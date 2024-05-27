package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/mrkovshik/yandex_diploma/api"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/service"
)

type restAPIServer struct {
	service api.Service
	storage service.Storage
	cfg     *config.Config
	logger  *zap.SugaredLogger
}

func NewRestAPIServer(service api.Service, storage service.Storage, cfg *config.Config, logger *zap.SugaredLogger) api.Server {
	return &restAPIServer{
		service: service,
		storage: storage,
		cfg:     cfg,
		logger:  logger,
	}
}
func (s *restAPIServer) RunServer(ctx context.Context) error {
	router := gin.Default()
	userSubRouter := router.Group("/api/user")
	{
		userSubRouter.POST("/register", s.RegisterHandler(ctx))
		userSubRouter.POST("/login", s.LoginHandler(ctx))
	}

	authGroup := userSubRouter.Group("")
	authGroup.Use(s.Auth(ctx))
	{
		authGroup.POST("/orders", s.UploadOrderHandler(ctx))
		authGroup.GET("/orders", s.GetOrders(ctx))
		authGroup.POST("/balance/withdraw", s.Withdraw(ctx))
		authGroup.GET("/balance", s.GetBalance(ctx))
		authGroup.GET("/withdrawals", s.ListWithdrawals(ctx))
	}

	return router.Run(s.cfg.RunAddress)
}
