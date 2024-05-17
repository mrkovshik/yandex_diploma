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

	"github.com/mrkovshik/yandex_diploma/internal/apperrors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func (s *restAPIServer) RegisterHandler(ctx context.Context) func(c *gin.Context) {
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

		token, err := s.service.Register(ctx, user.Login, user.Password)
		if err != nil {
			if errors.Is(err, apperrors.ErrUserAlreadyExists) {
				s.logger.Error("Register: ", err)
				c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			s.logger.Error("Register: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Header("Authorization", token)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "registration successful"})
	}
}

func (s *restAPIServer) LoginHandler(ctx context.Context) func(c *gin.Context) {
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
			if errors.Is(err, apperrors.ErrInvalidPassword) || errors.Is(err, sql.ErrNoRows) {
				s.logger.Error("Login: ", err)
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
				return
			}

			s.logger.Error("Login: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Header("Authorization", token)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "login successful"})
	}
}

func (s *restAPIServer) UploadOrderHandler(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		orderNumber, err := getOrderNumberFromContext(c)
		if err != nil {
			s.logger.Errorf("getOrderNumberFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid order number"})
			return
		}
		userID, err := getUserIDFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIDFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		exist, err := s.service.UploadOrder(ctx, orderNumber, userID)
		if err != nil {
			if errors.Is(err, apperrors.ErrOrderIsUploadedByAnotherUser) {
				s.logger.Error("UploadOrder", err)
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			s.logger.Error("UploadOrder", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if exist {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "order is already uploaded"})
			c.Abort()
			return
		}
		c.IndentedJSON(http.StatusAccepted, gin.H{"message": "order successfully uploaded"})
	}
}

func (s *restAPIServer) GetOrders(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIDFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		orders, err := s.service.GetUserOrders(ctx, userID)
		if err != nil {
			s.logger.Error("GetOrdersByUserID", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "no orders found"})
			c.Abort()
			return
		}
		c.IndentedJSON(http.StatusOK, orders)
	}
}

func (s *restAPIServer) Withdraw(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIDFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		var withdrawRequest model.Withdrawal
		if err := c.BindJSON(&withdrawRequest); err != nil {
			s.logger.Error("BindJSON", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := validate.Var(withdrawRequest.Amount, "required,min=1"); err != nil {
			s.logger.Error("validate Sum: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		orderNumberInt, err := strconv.ParseUint(withdrawRequest.OrderNumber, 10, 64)
		if err != nil {
			s.logger.Error("ParseUint: ", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := validate.Var(orderNumberInt, "required,luhn_checksum"); err != nil {
			s.logger.Error("validate OrderNumber: ", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		withdrawal := model.Withdrawal{
			Amount:      withdrawRequest.Amount,
			OrderNumber: withdrawRequest.OrderNumber,
			UserID:      userID,
		}
		if err := s.service.Withdraw(ctx, withdrawal); err != nil {
			if errors.Is(err, apperrors.ErrNotEnoughFunds) {
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

func (s *restAPIServer) GetBalance(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIDFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		balance, err1 := s.service.GetBalance(ctx, userID)
		if err1 != nil {
			s.logger.Errorf("GetBalance: %v", err1)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.IndentedJSON(http.StatusOK, balance)
	}
}

func (s *restAPIServer) ListWithdrawals(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := getUserIDFromContext(c)
		if err != nil {
			s.logger.Errorf("getUserIDFromContext: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid order number"})
			return
		}
		withdrawals, err1 := s.service.ListUserWithdrawals(ctx, userID)
		if err1 != nil {
			s.logger.Errorf("GetBalance: %v", err1)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if len(withdrawals) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.IndentedJSON(http.StatusOK, withdrawals)
	}
}
func getOrderNumberFromContext(c *gin.Context) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(c.Request.Body); err != nil {
		return "", err
	}
	//number, err1 := strconv.ParseUint(buf.String(), 10, 64)
	number := buf.String()
	err2 := validate.Var(number, "required,luhn_checksum")
	if err2 != nil {
		return "", err2
	}
	return number, nil
}

func getUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exist := c.Get("userID")
	if !exist {
		return 0, errors.New("no userID")
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		return 0, errors.New("error casting userID")
	}
	return userIDUint, nil
}
