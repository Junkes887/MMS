package repository

import (
	"github.com/Junkes887/MMS/internal/database/entity"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

func (r *Repository) FindMSS(pair string, from int64, to int64) ([]entity.MMSEntity, error) {
	var mmsEntity []entity.MMSEntity

	if err := r.CFG.DB.Find(&mmsEntity, "pair = ? AND timestamp BETWEEN ? AND ?", pair, from, to).Error; err != nil {
		return nil, err
	}

	r.Logger.Info("Consulta de MMS realizada com sucesso",
		zap.String("pair", pair),
		zap.Int64("from", from),
		zap.Int64("to", to),
		zap.Int("quantidade", len(mmsEntity)),
	)

	return mmsEntity, nil
}

func (r *Repository) SaveMSS(mmsEntity entity.MMSEntity) error {
	if err := r.CFG.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&mmsEntity).Error; err != nil {
		return err
	}
	r.Logger.Info("MMS salvo com sucesso",
		zap.String("pair", mmsEntity.Pair),
		zap.Int64("timestamp", mmsEntity.Timestamp),
	)
	return nil
}

func (r *Repository) BuscarDiasFaltantes(pair string, from, to int64) []int64 {
	var existentes []int64
	err := r.CFG.DB.Model(&entity.MMSEntity{}).
		Select("timestamp").
		Where("pair = ? AND timestamp BETWEEN ? AND ?", pair, from, to).
		Find(&existentes).Error

	if err != nil {
		r.Logger.Error("Erro ao buscar timestamps existentes", zap.Error(err))
		return nil
	}

	return existentes
}
