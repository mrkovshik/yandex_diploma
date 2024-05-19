package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/mrkovshik/yandex_diploma/internal/service/accrual"
	accrual2 "github.com/mrkovshik/yandex_diploma/mocks/accrual"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mrkovshik/yandex_diploma/internal/apperrors"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	"github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
	mock_service "github.com/mrkovshik/yandex_diploma/mocks"
)

const (
	UserLoginNotExist  = "none"
	orderExistingUser1 = "123456789049"
	orderExistingUser2 = "123456789015"
	orderNotExisting   = "123456789007"
	UserLogin1         = "JohnDow"
	UserPass1          = "qwerty"
	userHashedPass1    = "$2a$10$XVc79vBoRda4wdsx/uqMd.obXNtIbOvGttqUsgfBC4YfvuoD0fvrG"
	UserID1            = uint(123)
	UserID2            = uint(1232)
	withdrawalSumUser1 = 555.55
	balanceUser1       = float64(1000)
)

type GetOrderResp struct {
	ID          uint             `db:"id" json:"-"`
	OrderNumber string           `db:"order_number" json:"number"`
	UserID      uint             `db:"user_id" json:"-"`
	Status      model.OrderState `db:"status" json:"status"`
	UploadedAt  time.Time        `db:"uploaded_at" json:"uploaded_at"`
	Accrual     int              `db:"accrual" json:"accrual,omitempty"`
}

var (
	withdrawalUser1 = model.Withdrawal{
		ID:          324324,
		Amount:      300,
		ProcessedAt: time.Now(),
		OrderNumber: fmt.Sprint(orderExistingUser1),
		UserID:      UserID1,
	}
	withdrawalUser1a = model.Withdrawal{

		ID:          3433,
		Amount:      500,
		ProcessedAt: time.Now(),
		OrderNumber: fmt.Sprint(orderExistingUser1),
		UserID:      UserID1,
	}
)

func Test_restAPIServer_RunServer(t *testing.T) {
	var authToken string
	loggerConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "json", // You can use "console" for a more readable format
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "", // Disable caller key to remove caller information
			MessageKey:     "message",
			StacktraceKey:  "", // Disable stacktrace key to remove stack traces
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build the logger
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	ctx := context.Background()
	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("config.GetConfigs", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := defineStorage(ctx, ctrl)
	accrualService := accrual.NewAccrualService(cfg.AccrualSystemAddress)
	service := loyalty.NewBasicService(mockStorage, accrualService, cfg, sugar)
	srv := NewRestAPIServer(service, mockStorage, cfg, sugar)
	go func() {
		if err := srv.RunServer(ctx); err != nil {
			return
		}
	}()

	go accrual2.Run(cfg)

	time.Sleep(2 * time.Second)
	t.Run("register", func(t *testing.T) {
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLoginNotExist, UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/register", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		authToken = resp.Header().Get("Authorization")
		assert.NotEqual(t, 0, len(authToken))

		//User already exist
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLogin1, UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/register", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode())

		// No login
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, "", UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/register", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())

		// No password
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLogin1, "")).
			Post(fmt.Sprintf("http://%v/api/user/register", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("login", func(t *testing.T) {
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLogin1, UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/login", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		authToken = resp.Header().Get("Authorization")
		assert.NotEqual(t, 0, len(authToken))

		//User is not exist
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLoginNotExist, UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/login", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		// No login
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, "", UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/login", cfg.RunAddress))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("upload_order", func(t *testing.T) {
		url := fmt.Sprintf("http://%v/api/user/orders", cfg.RunAddress)
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			SetBody(fmt.Sprint(orderNotExisting)).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode())

		//Already uploaded
		resp1, err1 := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			SetBody(fmt.Sprint(orderExistingUser1)).
			Post(url)
		assert.NoError(t, err1)
		assert.Equal(t, http.StatusOK, resp1.StatusCode())

		//Already uploaded by another user
		resp2, err2 := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			SetBody(fmt.Sprint(orderExistingUser2)).
			Post(url)
		assert.NoError(t, err2)
		assert.Equal(t, http.StatusConflict, resp2.StatusCode())
	})

	t.Run("get_orders", func(t *testing.T) {
		url := fmt.Sprintf("http://%v/api/user/orders", cfg.RunAddress)
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		orders := make([]GetOrderResp, 3)
		err1 := json.Unmarshal(resp.Body(), &orders)
		assert.NoError(t, err1)
		assert.Equal(t, orders[0].OrderNumber, fmt.Sprint(orderExistingUser2))
		assert.Equal(t, orders[1].Status, model.OrderStateProcessing)
		assert.Equal(t, orders[2].Accrual, 1000)

		//No data
		resp2, err2 := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			Get(url)
		assert.NoError(t, err2)
		assert.Equal(t, http.StatusNoContent, resp2.StatusCode())
	})

	t.Run("get_balance", func(t *testing.T) {
		url := fmt.Sprintf("http://%v/api/user/balance", cfg.RunAddress)
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		balanceResp := model.GetBalanceResponse{}
		err1 := json.Unmarshal(resp.Body(), &balanceResp)
		assert.NoError(t, err1)
		assert.Equal(t, balanceResp.Balance, balanceUser1)
		assert.Equal(t, balanceResp.Withdrawn, withdrawalSumUser1)

	})

	t.Run("get_withdrawals", func(t *testing.T) {
		url := fmt.Sprintf("http://%v/api/user/withdrawals", cfg.RunAddress)
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", authToken).
			Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		withdrawals := make([]model.Withdrawal, 2)
		err1 := json.Unmarshal(resp.Body(), &withdrawals)
		assert.NoError(t, err1)
		assert.Equal(t, withdrawals[0].Amount, withdrawalUser1.Amount)
		assert.Equal(t, withdrawals[1].OrderNumber, withdrawalUser1a.OrderNumber)

	})

}

func defineStorage(ctx context.Context, ctrl *gomock.Controller) *mock_service.MockStorage {
	storage := mock_service.NewMockStorage(ctrl)
	storage.EXPECT().GetUserByLogin(ctx, UserLogin1).Return(model.User{
		ID:        UserID1,
		Login:     UserLogin1,
		Password:  userHashedPass1,
		Balance:   balanceUser1,
		CreatedAt: time.Now(),
	}, nil).AnyTimes()
	storage.EXPECT().GetUserByLogin(ctx, UserLoginNotExist).Return(model.User{}, sql.ErrNoRows).AnyTimes()
	storage.EXPECT().GetUserByID(ctx, UserID1).Return(model.User{
		ID:        UserID1,
		Login:     UserLogin1,
		Password:  userHashedPass1,
		Balance:   balanceUser1,
		CreatedAt: time.Now(),
	}, nil).AnyTimes()

	storage.EXPECT().AddUser(ctx, UserLoginNotExist, gomock.Any()).Return(UserID1, nil).AnyTimes()
	storage.EXPECT().AddUser(ctx, UserLogin1, gomock.Any()).Return(uint(0), apperrors.ErrUserAlreadyExists).AnyTimes()

	storage.EXPECT().GetOrderByNumber(ctx, orderNotExisting).Return(model.Order{}, sql.ErrNoRows).AnyTimes()
	storage.EXPECT().GetOrderByNumber(ctx, orderExistingUser1).Return(model.Order{
		ID:          878,
		OrderNumber: orderExistingUser1,
		UserID:      UserID1,
		Status:      model.OrderStateProcessed,
		UploadedAt:  time.Now(),
		Accrual:     1000,
	}, nil).AnyTimes()
	storage.EXPECT().GetOrderByNumber(ctx, orderExistingUser2).Return(model.Order{
		ID:          878,
		OrderNumber: orderExistingUser2,
		UserID:      UserID2,
		Status:      model.OrderStateProcessed,
		UploadedAt:  time.Now(),
		Accrual:     1000,
	}, nil).AnyTimes()
	gomock.InOrder(

		storage.EXPECT().GetOrdersByUserID(ctx, UserID1).Return([]model.Order{
			{
				ID:          878,
				OrderNumber: orderExistingUser2,
				UserID:      UserID2,
				Status:      model.OrderStateInvalid,
				UploadedAt:  time.Now(),
			},
			{
				ID:          234,
				OrderNumber: orderExistingUser1,
				UserID:      UserID2,
				Status:      model.OrderStateProcessing,
				UploadedAt:  time.Now(),
			},
			{
				ID:          8543,
				OrderNumber: orderNotExisting,
				UserID:      UserID2,
				Status:      model.OrderStateProcessed,
				UploadedAt:  time.Now(),
				Accrual:     1000,
			},
		}, nil),
		storage.EXPECT().GetOrdersByUserID(ctx, UserID1).Return([]model.Order{}, nil).AnyTimes(),
	)
	storage.EXPECT().UploadOrder(ctx, UserID1, orderNotExisting).Return(nil).AnyTimes().AnyTimes()

	storage.EXPECT().GetWithdrawalsSumByUserID(ctx, UserID1).Return(withdrawalSumUser1, nil).AnyTimes()

	storage.EXPECT().GetWithdrawalsByUserID(ctx, UserID1).Return([]model.Withdrawal{
		withdrawalUser1,
		withdrawalUser1a,
	}, nil).AnyTimes()

	return storage
}
