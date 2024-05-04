package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/internal/service"
)

func AddUserHandler(s service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var addUserReq addUserRequest
		if err := c.BindJSON(&addUserReq); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if err := s.AddUser(addUserReq.Login, addUserReq.Password); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}
}
