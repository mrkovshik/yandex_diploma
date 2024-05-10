package counting

type Service interface {
	GetOrderScore(orderNumber int) (Response, error)
}
