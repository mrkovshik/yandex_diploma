package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/server"
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
	srv := server.NewServer(sugar)
	run(srv, ctx, cfg)
}

func run(s *server.Server, ctx context.Context, cfg *config.Config) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
