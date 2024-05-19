package accrual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/mrkovshik/yandex_diploma/internal/apperrors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
	"github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
)

type service struct {
	address string
}

func NewAccrualService(address string) loyalty.AccrualService {
	return service{
		address: address,
	}
}

func (s service) GetOrderAccrual(orderNumber string) (model.AccrualResponse, error) {
	var orderResponse model.AccrualResponse
	serviceURL := fmt.Sprintf("%v/api/orders/%v", s.address, orderNumber)
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(60 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 60 * time.Second, apperrors.ErrTooManyRetrials
		})
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	resp, err := client.R().Get(serviceURL)
	if err != nil {
		return model.AccrualResponse{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNoContent {
			return model.AccrualResponse{}, apperrors.ErrNoSuchOrder
		}
		return model.AccrualResponse{}, fmt.Errorf("status code: %v", resp.StatusCode())
	}
	if err := json.Unmarshal(resp.Body(), &orderResponse); err != nil {
		return model.AccrualResponse{}, err
	}
	return orderResponse, nil
}
