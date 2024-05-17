package accrual

type Service interface {
	GetOrderAccrual(orderNumber string) (Response, error)
}
