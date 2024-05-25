package main

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mrkovshik/yandex_diploma/api/rest"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/service/accrual"
	"github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
	"github.com/mrkovshik/yandex_diploma/internal/storage/postgres"
)

const accrualInterval = 10 * time.Second //TODO: move to config

func main() {
	loggerConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "",
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	ctx := context.Background()
	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("config.GetConfigs", err)
	}
	db, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	if err != nil {
		sugar.Fatal("sql.Open", err)
	}
	db.MustExec(postgres.Schema)
	accrualService := accrual.NewAccrualService(cfg.AccrualSystemAddress)
	storage := postgres.NewStorage(db)
	service := loyalty.NewBasicService(storage, accrualService, cfg, sugar)

	srv := rest.NewRestAPIServer(service, storage, cfg, sugar)

	accrualTicker := time.NewTicker(accrualInterval)
	go func() {
		for range accrualTicker.C {
			if err := service.UpdatePendingOrders(ctx); err != nil {
				sugar.Errorf("UpdatePendingOrders: %e", err)
			}
		}
	}()

	if err := srv.RunServer(ctx); err != nil {
		sugar.Fatal(err)
	}
	sugar.Fatal("RunServer", err)
}
