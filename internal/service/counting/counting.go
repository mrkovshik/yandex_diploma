package counting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type (
	service struct {
		address string
	}
	Response struct {
		Order   int                 `json:"order"`
		Status  model.CountingState `json:"status"`
		Accrual int                 `json:"accrual"`
	}
)

func NewCountingService(address string) Service {
	return service{
		address: address,
	}
}

func (s service) GetOrderAccrual(orderNumber uint) (Response, error) {
	var orderResponse Response
	serviceURL := fmt.Sprintf("http://%v/api/orders/%v", s.address, orderNumber)
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, app_errors.ErrTooManyRetrials
		})
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	resp, err := client.R().Get(serviceURL)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNoContent {
			return Response{}, app_errors.ErrNoSuchOrder
		}
		return Response{}, app_errors.ErrInvalidResponseCode
	}
	if err := json.Unmarshal(resp.Body(), &orderResponse); err != nil {
		return Response{}, err
	}
	return orderResponse, nil
}
