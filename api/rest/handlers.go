package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/internal/service"
)

func (s restApiServer) RegisterHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var addUserReq addUserRequest
		basicService := service.NewBasicService(s.db, s.cfg, s.logger)
		if err := c.BindJSON(&addUserReq); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		if err := basicService.Register(ctx, addUserReq.Login, addUserReq.Password); err != nil {
			s.logger.Error("AddUser", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "user successfully registered"})
	}
}
