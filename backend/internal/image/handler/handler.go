package handler

import imageservice "github.com/fiqryomaratala/image-processing-service/backend/internal/image/service"

type Handler struct {
	service imageservice.Service
}

func NewHandler(service imageservice.Service) *Handler {
	return &Handler{
		service: service,
	}
}
