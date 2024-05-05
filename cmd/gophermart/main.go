package main

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/mrkovshik/yandex_diploma/api/rest"
	"github.com/mrkovshik/yandex_diploma/internal/config"
)

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
	order_number integer NOT NULL,
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

	srv := rest.NewRestApiServer(db, cfg, sugar)

	if err := srv.RunServer(ctx); err != nil {
	}
	sugar.Fatal("RunServer", err)
}
