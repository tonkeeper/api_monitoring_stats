package controllers

import (
	"context"

	"api_monitoring_stats/controllers/oas"
)

var _ oas.Handler = (*Handler)(nil)

type Handler struct {
	oas.UnimplementedHandler // automatically implement all methods
}

func NewHandler() (*Handler, error) {
	return &Handler{}, nil
}

func (h Handler) PingReadyGet(ctx context.Context) (oas.PingReadyGetRes, error) {
	return &oas.PingReadyGetOK{}, nil
}

func (h Handler) PingReadyHead(ctx context.Context) (oas.PingReadyHeadRes, error) {
	return &oas.PingReadyHeadOK{}, nil
}
