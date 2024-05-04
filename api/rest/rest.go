package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/api"
	"github.com/mrkovshik/yandex_diploma/internal/service"
)

type restApiServer struct {
	service service.Service
}

func NewRestApiServer(service service.Service) api.Server {
	return restApiServer{
		service: service,
	}
}
func (s restApiServer) RunServer(ctx context.Context) error {
	router := gin.Default()
	userSubRouter := router.Group("/api/user")
	userSubRouter.POST("/register", AddUserHandler(s.service))
	return router.Run("localhost:8080")
}
