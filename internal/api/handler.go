package api

import (
	"github.com/Junkes887/MMS/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	Service service.ServiceInterface
	Logger  *zap.Logger
}

func NewHandler(service service.ServiceInterface, logger *zap.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}
