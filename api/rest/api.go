package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mrkovshik/yandex_diploma/api"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"go.uber.org/zap"
)

type restApiServer struct {
	db     *sqlx.DB
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewRestApiServer(db *sqlx.DB, cfg *config.Config, logger *zap.SugaredLogger) api.Server {
	return restApiServer{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}
func (s restApiServer) RunServer(ctx context.Context) error {
	router := gin.Default()
	userSubRouter := router.Group("/api/user")
	userSubRouter.POST("/register", s.RegisterHandler(ctx))
	userSubRouter.POST("/login", s.LoginHandler(ctx))
	return router.Run("localhost:8080")
}
