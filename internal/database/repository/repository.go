package repository

import (
	"github.com/Junkes887/MMS/internal/database"
	"github.com/Junkes887/MMS/internal/database/entity"
	"go.uber.org/zap"
)

type RepositoryInterface interface {
	FindMSS(pair string, from, to int64) ([]entity.MMSEntity, error)
	SaveMSS(mmsEntity entity.MMSEntity) error
	BuscarDiasFaltantes(pair string, from, to int64) []int64
}

type Repository struct {
	CFG    *database.DBConnection
	Logger *zap.Logger
}

func NewRepository(cfg *database.DBConnection, logger *zap.Logger) *Repository {
	return &Repository{
		CFG:    cfg,
		Logger: logger,
	}
}
