package fetcher

import (
	"github.com/Junkes887/MMS/internal/database/repository"
	"go.uber.org/zap"
)

type Fetcher struct {
	Repository repository.RepositoryInterface
	Logger     *zap.Logger
}

func NewFetcher(repository repository.RepositoryInterface, logger *zap.Logger) *Fetcher {
	return &Fetcher{
		Repository: repository,
		Logger:     logger,
	}
}
