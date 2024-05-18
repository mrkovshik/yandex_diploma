package loyalty

import (
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

type AccrualService interface {
	GetOrderAccrual(orderNumber string) (model.AccrualResponse, error)
}
