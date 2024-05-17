package api

import "github.com/mrkovshik/yandex_diploma/internal/service/accrual"

type AccrualService interface {
	GetOrderAccrual(orderNumber string) (accrual.Response, error)
}
