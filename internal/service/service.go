package service

type Service interface {
	AddUser(login, password string) error
}
