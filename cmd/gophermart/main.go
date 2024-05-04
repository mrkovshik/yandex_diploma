package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

	srv := server.NewServer(sugar)
	run(srv, ctx)
}

func run(s *server.Server, ctx context.Context) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
