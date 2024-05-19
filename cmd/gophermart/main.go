package main

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mrkovshik/yandex_diploma/internal/service/accrual"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mrkovshik/yandex_diploma/api/rest"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
	"github.com/mrkovshik/yandex_diploma/internal/storage/postgres"
)

const accrualInterval = 10 * time.Second //TODO: move to config
var schema = `
CREATE TABLE IF NOT EXISTS users (
	id serial4 NOT NULL,
	login varchar NOT NULL,
	"password" varchar NOT NULL,
	created_at timestamptz NOT NULL,
	balance float4 DEFAULT 0 NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS orders (
	id serial4 NOT NULL,
	order_number varchar NOT NULL,
	user_id int4 NOT NULL,
	uploaded_at timestamptz NOT NULL,
	status varchar DEFAULT 'NEW'::character varying NOT NULL,
	accrual float4 DEFAULT 0 NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id),
	CONSTRAINT orders_unique UNIQUE (order_number)
);                                  
    CREATE TABLE IF NOT EXISTS withdrawals (
		id serial4 NOT NULL,
	amount float4 NOT NULL,
	processed_at timestamptz NOT NULL,
	order_number varchar NOT NULL,
	user_id int4 NOT NULL,
	CONSTRAINT withdrawals_pk PRIMARY KEY (id)
)`

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
	db.MustExec(schema)
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
