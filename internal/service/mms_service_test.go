package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Junkes887/MMS/internal/database/entity"
	"github.com/Junkes887/MMS/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMMSBuscar_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := &Service{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	from := time.Now().AddDate(0, 0, -5).Unix()
	to := time.Now().Unix()

	mockData := []entity.MMSEntity{
		{Pair: "BRLBTC", Timestamp: from, MMS20: 50000.0},
		{Pair: "BRLBTC", Timestamp: from + 86400, MMS20: 51000.0},
		{Pair: "BRLBTC", Timestamp: from + 2*86400, MMS20: 52000.0},
	}

	mockRepo.On("FindMSS", "BRLBTC", from, to).Return(mockData, nil)

	result, err := service.MMSBuscar("BRLBTC", 20, from, to)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, mockData[0].MMS20, result[0].MMS)
	assert.Equal(t, mockData[1].MMS20, result[1].MMS)
}

func TestMMSBuscar_ErrorPair(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := &Service{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	from := time.Now().AddDate(0, 0, -5).Unix()
	to := time.Now().Unix()

	result, err := service.MMSBuscar("PAIR ERRADO", 2, from, to)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "pair informado é inválido", err.Error())
}

func TestMMSBuscar_ErrorFromDatabase(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := &Service{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	from := time.Now().AddDate(0, 0, -5).Unix()
	to := time.Now().Unix()

	mockRepo.On("FindMSS", "BRLBTC", from, to).Return([]entity.MMSEntity{}, errors.New("db error"))

	result, err := service.MMSBuscar("BRLBTC", 2, from, to)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "erro ao consultar dados no banco", err.Error())
}

func TestMMSBuscar_FromTooOld(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := &Service{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	from := time.Now().AddDate(0, 0, -366).Unix()
	to := time.Now().Unix()

	result, err := service.MMSBuscar("BRLBTC", 2, from, to)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "from não pode ser anterior a 365 dias", err.Error())
}
