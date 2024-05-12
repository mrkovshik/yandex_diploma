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
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	"github.com/mrkovshik/yandex_diploma/internal/service/accrual/mock"
	mock_loyalty "github.com/mrkovshik/yandex_diploma/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type tokenResp struct {
	Token string
}

const (
	UserLoginNotExist        = "none"
	orderExisting            = uint(1234567890)
	orderNotExisting         = uint(123456789007)
	NumberForTooManyRequests = uint(2468013579)
	UserLogin1               = "JohnDow"
	UserPass1                = "qwerty"
	userHashedPass1          = "$2a$10$XVc79vBoRda4wdsx/uqMd.obXNtIbOvGttqUsgfBC4YfvuoD0fvrG"
	UserId1                  = uint(123)
	UserIdNotExist           = uint(456)
)

func Test_restApiServer_RunServer(t *testing.T) {
	var authToken tokenResp
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()

	ctx := context.Background()
	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("config.GetConfigs", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := defineStorage(ctx, ctrl)

	srv := NewRestApiServer(mockStorage, cfg, sugar)
	go func() {
		if err := srv.RunServer(ctx); err != nil {
			return
		}
	}()
	go mock.Run(cfg)

	t.Run("register", func(t *testing.T) {
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLoginNotExist, UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/register", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

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
		err1 := json.Unmarshal(resp.Body(), &authToken)
		assert.NoError(t, err1)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		//User is not exist
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, UserLoginNotExist, UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/login", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
		assert.NotEmpty(t, resp.Body())
		// No login
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"login":"%v", "password":"%v"}`, "", UserPass1)).
			Post(fmt.Sprintf("http://%v/api/user/login", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("upload_order", func(t *testing.T) {
		client := resty.New()

		//Normal flow
		resp, err := client.R().SetHeader("Content-Type", "text/plain").
			SetHeader("Authorization", fmt.Sprintf("Bearer %v", authToken.Token)).
			SetBody(fmt.Sprint(orderNotExisting)).
			Post(fmt.Sprintf("http://%v/api/user/orders", cfg.RunAddress))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	//	//Already exist
	//	resp, err := client.R().SetHeader("Content-Type", "text/plain").
	//		SetHeader("Authorization", fmt.Sprintf("Bearer %v", authToken.Token)).
	//		SetBody(fmt.Sprint(orderNotExisting)).
	//		Post(fmt.Sprintf("http://%v/api/user/orders", cfg.RunAddress))
	//	assert.NoError(t, err)
	//	assert.Equal(t, http.StatusOK, resp.StatusCode())
	//})

}

func defineStorage(ctx context.Context, ctrl *gomock.Controller) *mock_loyalty.MockStorage {
	storage := mock_loyalty.NewMockStorage(ctrl)
	storage.EXPECT().GetUserByLogin(ctx, UserLogin1).Return(model.User{
		ID:        UserId1,
		Login:     UserLogin1,
		Password:  userHashedPass1,
		Balance:   50,
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{},
	}, nil)
	storage.EXPECT().GetUserByLogin(ctx, UserLoginNotExist).Return(model.User{}, sql.ErrNoRows)
	//storage.EXPECT().GetUserByID(ctx, UserIdNotExist).Return(model.User{}, sql.ErrNoRows)
	storage.EXPECT().GetUserByID(ctx, UserId1).Return(model.User{
		ID:        UserId1,
		Login:     UserLogin1,
		Password:  userHashedPass1,
		Balance:   50,
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{},
	}, nil)

	storage.EXPECT().AddUser(ctx, UserLoginNotExist, gomock.Any()).Return(nil)
	storage.EXPECT().AddUser(ctx, UserLogin1, gomock.Any()).Return(app_errors.ErrUserAlreadyExists)

	storage.EXPECT().GetOrderByNumber(ctx, orderNotExisting).Return(model.Order{}, sql.ErrNoRows)

	storage.EXPECT().UploadOrder(ctx, UserId1, orderNotExisting).Return(nil)

	//storage.EXPECT().GetOrderByNumber(ctx, orderExisting).Return(model.Order{
	//	ID:          5,
	//	OrderNumber: orderExisting,
	//	UserId:      UserId1,
	//	Status:      model.OrderStateProcessed,
	//	UploadedAt:  time.Now(),
	//	UpdatedAt:   sql.NullTime{},
	//}, nil)
	return storage
}
