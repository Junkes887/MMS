package fetcher

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Junkes887/MMS/internal/constants"
	"github.com/Junkes887/MMS/internal/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestSeedData_Success(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	fetcher := &Fetcher{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	mockRepo.On("SaveMSS", mock.Anything).Return(nil)

	fetcher.SeedData("BTC-BRL", "BRLBTC", time.Now().AddDate(0, 0, -365).Unix(), time.Now().Unix())

	mockRepo.AssertExpectations(t)
}

func TestSeedData_FetchCandlesError(t *testing.T) {
	os.Setenv("DELAY_FETCHER_SECONDS", "0")
	defer os.Unsetenv("DELAY_FETCHER_SECONDS")

	mockRepo := new(mocks.MockRepository)
	fetcher := &Fetcher{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	originalFetchCandles := fetchCandles
	defer func() { fetchCandles = originalFetchCandles }()

	fetchCandles = func(symbol string, from, to int64) ([]float64, []int64, error) {
		return nil, nil, errors.New("erro ao buscar candles")
	}

	fetcher.SeedData("BTC-BRL", "BRLBTC", time.Now().AddDate(0, 0, -1).Unix(), time.Now().Unix())

	mockRepo.AssertNotCalled(t, "SaveMSS", mock.Anything)
}

func TestSeedData_SaveMSSError(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	fetcher := &Fetcher{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	mockRepo.On("SaveMSS", mock.Anything).Return(errors.New("erro ao salvar"))

	originalFetchCandles := fetchCandles
	defer func() { fetchCandles = originalFetchCandles }()

	fetchCandles = func(symbol string, from, to int64) ([]float64, []int64, error) {
		return []float64{100.0, 200.0, 300.0}, []int64{time.Now().Unix(), time.Now().Unix(), time.Now().Unix()}, nil
	}

	fetcher.SeedData("BTC-BRL", "BRLBTC", time.Now().AddDate(0, 0, -1).Unix(), time.Now().Unix())

	mockRepo.AssertExpectations(t)
}

func TestVerificarDadosFaltantes(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	fetcher := &Fetcher{
		Repository: mockRepo,
		Logger:     zap.NewNop(),
	}

	constants.SymbolPairMap = map[string]string{
		"BTC-BRL": "BRLBTC",
	}

	now := time.Now().Truncate(24 * time.Hour)
	mockRepo.On("BuscarDiasFaltantes", "BRLBTC", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
		Return([]int64{now.AddDate(0, 0, 0).Unix()})

	emailAlertaChamado := make(chan struct{}, 1)

	originalEnviarEmail := enviarEmailAlerta
	defer func() { enviarEmailAlerta = originalEnviarEmail }()

	enviarEmailAlerta = func(f *Fetcher, pair string, diasFaltando []int64) {
		assert.Equal(t, "BRLBTC", pair)
		assert.Equal(t, 365, len(diasFaltando))
		emailAlertaChamado <- struct{}{}
	}

	fetcher.VerificarDadosFaltantes()

	select {
	case <-emailAlertaChamado:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando email de alerta ser enviado")
	}

	mockRepo.AssertExpectations(t)
}
