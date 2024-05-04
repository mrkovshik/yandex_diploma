package server

import (
	"database/sql"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	db     *sql.DB
}

func NewServer(logger *zap.SugaredLogger) *Server {
	return &Server{
		logger: logger,
	}
}
