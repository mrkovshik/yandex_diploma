package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/internal/service"
	"go.uber.org/zap"

	"github.com/mrkovshik/yandex_diploma/api"
	"github.com/mrkovshik/yandex_diploma/internal/config"
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
	userSubRouter.POST("/register", s.RegisterHandler(ctx))
	userSubRouter.POST("/login", s.LoginHandler(ctx))
	userSubRouter.POST("/orders", s.Auth(ctx), s.UploadOrderHandler(ctx))
	userSubRouter.GET("/orders", s.Auth(ctx), s.GetOrders(ctx))
	userSubRouter.POST("/balance/withdraw", s.Auth(ctx), s.Withdraw(ctx))
	userSubRouter.GET("/balance", s.Auth(ctx), s.GetBalance(ctx))
	userSubRouter.GET("/withdrawals", s.Auth(ctx), s.ListWithdrawals(ctx))
	return router.Run("localhost:8080")
}
