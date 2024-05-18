package main

import (
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/mocks"
)

func main() {
	cfg, _ := config.GetConfigs()
	mock_api.Run(cfg)
}
