package mocks

import (
	"github.com/Junkes887/MMS/internal/database/entity"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindMSS(pair string, from, to int64) ([]entity.MMSEntity, error) {
	args := m.Called(pair, from, to)
	return args.Get(0).([]entity.MMSEntity), args.Error(1)
}

func (m *MockRepository) SaveMSS(mmsEntity entity.MMSEntity) error {
	args := m.Called(mmsEntity)
	return args.Error(0)
}

func (m *MockRepository) BuscarDiasFaltantes(pair string, from, to int64) []int64 {
	args := m.Called(pair, from, to)
	return args.Get(0).([]int64)
}
