package handlers

import "gpsd-api-gateway/internal/gateway/pkg/config"

type Handler struct {
	Config *config.Config
}

func NewHandler(cc *config.Config) *Handler {
	return &Handler{Config: cc}
}
