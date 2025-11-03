package handlers

import (
	"github.com/notLeoHirano/bartr/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}