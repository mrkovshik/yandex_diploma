package counting

type Service interface {
	GetOrderAccrual(orderNumber uint) (Response, error)
}
