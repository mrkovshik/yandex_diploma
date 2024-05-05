package service

import "context"

type Service interface {
	Register(ctx context.Context, login, password string) error
}
