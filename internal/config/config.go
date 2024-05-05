package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	TokenExp             int    `env:"TOKEN_EXP" envDefault:"3"`
	SecretKey            string `env:"SECRET_KEY" envDefault:"MyBaby'sGotASecret"`
}

type serverConfigBuilder struct {
	Config Config
}

func (c *serverConfigBuilder) withRunAddress(host string) *serverConfigBuilder {
	c.Config.RunAddress = host
	return c
}

func (c *serverConfigBuilder) withDatabaseURI(uri string) *serverConfigBuilder {
	c.Config.DatabaseURI = uri
	return c
}

func (c *serverConfigBuilder) withAccrualSystemAddress(dsn string) *serverConfigBuilder {
	c.Config.AccrualSystemAddress = dsn
	return c
}

func (c *serverConfigBuilder) fromFlags() *serverConfigBuilder {
	runAddress := flag.String("a", "localhost:8080", "server host and port")
	databaseURI := flag.String("d", "postgres://yandex:yandex@localhost:5432/yandex?sslmode=disable", "db URI")
	accrualSystemAddress := flag.String("r", "localhost:8081", "accrual system host and port")
	flag.Parse()

	if c.Config.RunAddress == "" {
		c.withRunAddress(*runAddress)
	}
	if c.Config.DatabaseURI == "" {
		c.withDatabaseURI(*databaseURI)
	}
	if c.Config.AccrualSystemAddress == "" {
		c.withAccrualSystemAddress(*accrualSystemAddress)
	}

	return c
}

func (c *serverConfigBuilder) fromEnv() *serverConfigBuilder {

	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	return c
}

func GetConfigs() (*Config, error) {
	var c serverConfigBuilder
	c.fromEnv().fromFlags()
	//TODO: add validation here
	return &c.Config, nil
}
