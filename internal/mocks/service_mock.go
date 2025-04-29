package mocks

import (
	"github.com/Junkes887/MMS/internal/service/dto"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) MMSBuscar(pair string, mmsRange int, from int64, to int64) ([]dto.MMSResponse, error) {
	args := m.Called(pair, mmsRange, from, to)
	return args.Get(0).([]dto.MMSResponse), args.Error(1)
}
