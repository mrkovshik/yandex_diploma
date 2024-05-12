package main

import (
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/service/accrual/mock"
)

func main() {
	cfg, _ := config.GetConfigs()
	mock.Run(cfg)
}
