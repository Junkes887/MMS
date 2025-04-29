package service

import (
	"github.com/Junkes887/MMS/internal/database/repository"
	"github.com/Junkes887/MMS/internal/service/dto"
	"go.uber.org/zap"
)

type ServiceInterface interface {
	MMSBuscar(pair string, mmsRange int, from int64, to int64) ([]dto.MMSResponse, error)
}

type Service struct {
	Repository repository.RepositoryInterface
	Logger     *zap.Logger
}

func NewService(repository repository.RepositoryInterface, logger *zap.Logger) *Service {
	return &Service{
		Repository: repository,
		Logger:     logger,
	}
}
