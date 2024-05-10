package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/api"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
	"go.uber.org/zap"
)

type restApiServer struct {
	userStorage  loyalty.UserStorage
	orderStorage loyalty.OrderStorage
	cfg          *config.Config
	logger       *zap.SugaredLogger
}

func NewRestApiServer(userStorage loyalty.UserStorage, orderStorage loyalty.OrderStorage, cfg *config.Config, logger *zap.SugaredLogger) api.Server {
	return restApiServer{
		userStorage:  userStorage,
		orderStorage: orderStorage,
		cfg:          cfg,
		logger:       logger,
	}
}
func (s restApiServer) RunServer(ctx context.Context) error {
	router := gin.Default()
	userSubRouter := router.Group("/api/user")
	userSubRouter.POST("/register", s.RegisterHandler(ctx))
	userSubRouter.POST("/login", s.LoginHandler(ctx))
	userSubRouter.POST("/orders", s.Auth(ctx), s.UploadOrderHandler(ctx))
	return router.Run("localhost:8080")
}
