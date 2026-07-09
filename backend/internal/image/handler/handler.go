package handler

import (
	imageservice "github.com/fiqryomaratala/image-processing-service/backend/internal/image/service"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"go.uber.org/zap"
)

type Handler struct {
	service imageservice.Service
	logger  *zap.Logger
}

func NewHandler(service imageservice.Service) *Handler {
	log := zap.NewNop()

	defer func() {
		_ = recover()
	}()

	log = logger.Get()

	return &Handler{
		service: service,
		logger:  log,
	}
}
