package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mrkovshik/yandex_diploma/internal/service/loyalty"
	"github.com/mrkovshik/yandex_diploma/internal/storage/postgres"
	"go.uber.org/zap"

	"github.com/mrkovshik/yandex_diploma/api/rest"
	"github.com/mrkovshik/yandex_diploma/internal/config"
)

const accrualInterval = 60 * time.Second //TODO: move to config
var schema = `
CREATE TABLE IF NOT EXISTS users (
	id serial NOT NULL,
	login varchar NOT NULL,
	"password" varchar NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NULL,
	CONSTRAINT users_pk PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS orders (
	id serial NOT NULL,
	order_number bigint NOT NULL,
	user_id integer NOT NULL,
	uploaded_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NULL,
	status varchar NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id),
	CONSTRAINT orders_unique UNIQUE (order_number),
	CONSTRAINT orders_users_fk FOREIGN KEY (user_id) REFERENCES users(id)
)`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

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
	db, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	if err != nil {
		sugar.Fatal("sql.Open", err)
	}
	db.MustExec(schema)
	storage := postgres.NewStorage(db)
	service := loyalty.NewBasicService(storage, cfg, sugar)
	srv := rest.NewRestApiServer(service, storage, cfg, sugar)

	accrualTicker := time.NewTicker(accrualInterval)
	go func() {
		for range accrualTicker.C {
			if err := service.UpdatePendingOrders(ctx); err != nil {
				sugar.Errorf("UpdatePendingOrders: %e", err)
			}
		}
	}()

	if err := srv.RunServer(ctx); err != nil {
	}
	sugar.Fatal("RunServer", err)
}
