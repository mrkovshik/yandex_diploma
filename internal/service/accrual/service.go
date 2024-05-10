package accrual

type Service interface {
	GetOrderAccrual(orderNumber uint) (Response, error)
}
