package main

import (
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/mocks/accrual"
)

func main() {
	cfg, _ := config.GetConfigs()
	accrual.Run(cfg)
}
