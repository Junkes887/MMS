package service

import (
	"errors"
	"time"

	"github.com/Junkes887/MMS/internal/constants"
	"github.com/Junkes887/MMS/internal/service/dto"
	"github.com/Junkes887/MMS/pkg/mapper"
	"go.uber.org/zap"
)

func (s *Service) MMSBuscar(pair string, mmsRange int, from int64, to int64) ([]dto.MMSResponse, error) {
	found := false
	for _, v := range constants.SymbolPairMap {
		if v == pair {
			found = true
			break
		}
	}
	if !found {
		msg := "pair informado é inválido"
		s.Logger.Warn(msg, zap.String("pair", pair))
		return nil, errors.New(msg)
	}

	limitTimestamp := time.Now().AddDate(0, 0, -365).Unix()
	if from < limitTimestamp {
		msg := "from não pode ser anterior a 365 dias"
		s.Logger.Warn(msg, zap.Int64("from", from))
		return nil, errors.New(msg)
	}

	list, err := s.Repository.FindMSS(pair, from, to)
	if err != nil {
		msg := "erro ao consultar dados no banco"
		s.Logger.Error(msg, zap.Error(err))
		return nil, errors.New(msg)
	}

	s.Logger.Info("Consulta no banco realizada com sucesso",
		zap.String("pair", pair),
		zap.Int64("from", from),
		zap.Int64("to", to),
		zap.Int("resultados", len(list)),
	)

	return mapper.MMSEntityToMSSResponse(list, mmsRange), nil
}
