package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/service"
)

func (s restApiServer) RegisterHandler(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		var addUserReq addUserRequest
		basicService := service.NewBasicService(s.db, s.cfg, s.logger)
		if err := c.BindJSON(&addUserReq); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := basicService.Register(ctx, addUserReq.Login, addUserReq.Password); err != nil {
			if errors.Is(err, app_errors.ErrUserAlreadyExists) {
				s.logger.Error("Register: ", err)
				c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
			}
			s.logger.Error("Register: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "user successfully registered"})
	}
}
