package rest

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func (s *restApiServer) RegisterHandler(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User
		if err := c.BindJSON(&user); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := validate.Struct(user); err != nil {
			s.logger.Error("validate.Struct", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := s.service.Register(ctx, user.Login, user.Password); err != nil {
			if errors.Is(err, app_errors.ErrUserAlreadyExists) {
				s.logger.Error("Register: ", err)
				c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			s.logger.Error("Register: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "user successfully registered"})
	}
}

func (s *restApiServer) LoginHandler(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User

		if err := c.BindJSON(&user); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := validate.Struct(user); err != nil {
			s.logger.Error("validate.Struct", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := s.service.Login(ctx, user.Login, user.Password)
		if err != nil {
			if errors.Is(err, app_errors.ErrInvalidPassword) || errors.Is(err, sql.ErrNoRows) {
				s.logger.Error("Register: ", err)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}

			s.logger.Error("Register: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"token": token})
	}
}

func (s *restApiServer) UploadOrderHandler(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		orderNumber, err := getOrderNumberFromContext(c)
		if err != nil {
			s.logger.Errorf("getOrderNumberFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid order number"})
			return
		}
		userId, err := getUserIdFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIdFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		exist, err := s.service.UploadOrder(ctx, orderNumber, userId)
		if err != nil {
			if errors.Is(err, app_errors.ErrOrderIsUploadedByAnotherUser) {
				s.logger.Error("UploadOrder", err)
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			s.logger.Error("UploadOrder", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if exist {
			c.IndentedJSON(http.StatusAccepted, gin.H{"message": "order is already uploaded"})
			c.Abort()
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "order successfully uploaded"})
	}
}

func (s *restApiServer) GetOrders(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := getUserIdFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIdFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		orders, err := s.service.GetUserOrders(ctx, userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.IndentedJSON(http.StatusNoContent, gin.H{"message": "no orders found"})
				c.Abort()
				return
			}
			s.logger.Error("GetOrdersByUserID", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, orders)
	}
}

func (s *restApiServer) Withdraw(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := getUserIdFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIdFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		var withdrawRequest model.WithdrawRequest
		if err := c.BindJSON(&withdrawRequest); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := validate.Var(withdrawRequest.Sum, "required,min=1"); err != nil {
			s.logger.Error("validate Sum: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		orderNumber, err := strconv.ParseUint(withdrawRequest.Order, 10, 64)
		if err != nil {
			s.logger.Error("ParseUint: ", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := validate.Var(orderNumber, "required,luhn_checksum"); err != nil {
			s.logger.Error("validate OrderNumber: ", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		withdrawal := model.Withdrawal{
			Amount:      withdrawRequest.Sum,
			OrderNumber: uint(orderNumber),
			UserId:      userId,
		}
		if err := s.service.Withdraw(ctx, withdrawal); err != nil {
			if errors.Is(err, app_errors.ErrNotEnoughFunds) {
				s.logger.Error("Withdraw", err)
				c.IndentedJSON(http.StatusPaymentRequired, gin.H{"message": err.Error()})
				c.Abort()
				return
			}
			s.logger.Error("Withdraw", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "withdrawal successfully processed"})
	}
}

func getOrderNumberFromContext(c *gin.Context) (uint, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(c.Request.Body); err != nil {
		return 0, err
	}
	number, err1 := strconv.ParseUint(buf.String(), 10, 64)
	if err1 != nil {
		return 0, err1
	}
	err2 := validate.Var(number, "required,luhn_checksum")
	if err2 != nil {
		return 0, err2
	}
	return uint(number), nil
}

func getUserIdFromContext(c *gin.Context) (uint, error) {
	userId, exist := c.Get("userId")
	if !exist {
		return 0, errors.New("no userId")
	}
	userIdUint, ok := userId.(uint)
	if !ok {
		return 0, errors.New("error casting userId")
	}
	return userIdUint, nil
}
