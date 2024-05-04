package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mrkovshik/yandex_diploma/api/rest"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()
	ctx := context.Background()
	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("config.GetConfigs", err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		sugar.Fatal("sql.Open", err)
	}
	svc := service.NewBasicService(db, cfg, sugar)
	srv := rest.NewRestApiServer(svc)

	if err := srv.RunServer(ctx); err != nil {
	}
	sugar.Fatal("RunServer", err)
}
