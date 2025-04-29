package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Junkes887/MMS/internal/mocks"
	"github.com/Junkes887/MMS/internal/service/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:pair/mms", handler.MMSBuscar)
	return router
}

func TestMMSBuscar_Success(t *testing.T) {
	mockService := new(mocks.MockService)
	logger := zap.NewNop()
	handler := &Handler{Service: mockService, Logger: logger}

	router := setupRouter(handler)

	from := int64(1713744000)
	to := int64(1713830400)

	expectedResp := []dto.MMSResponse{
		{Timestamp: from, MMS: 50000.0},
		{Timestamp: to, MMS: 51000.0},
	}

	mockService.On("MMSBuscar", "BRLBTC", 20, from, to).Return(expectedResp, nil)

	req, _ := http.NewRequest("GET", "/BRLBTC/mms?range=20&from="+
		strconv.FormatInt(from, 10)+"&to="+strconv.FormatInt(to, 10), nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "50000")
	assert.Contains(t, resp.Body.String(), "51000")
}

func TestMMSBuscar_BadRequestRange(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := &Handler{Service: mockService}

	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/BRLBTC/mms?range=15&from=1713744000", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "range inválido")
}

func TestMMSBuscar_BadRequestFrom(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := &Handler{Service: mockService}

	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/BRLBTC/mms?range=20&from=abc", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "from inválido")
}

func TestMMSBuscar_BadRequestTo(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := &Handler{Service: mockService}
	from := int64(1713744000)

	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/BRLBTC/mms?range=20&from="+
		strconv.FormatInt(from, 10)+"&to=abc", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "to inválido")
}

func TestMMSBuscar_ErrorFromService(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := &Handler{Service: mockService}

	router := setupRouter(handler)

	from := time.Now().AddDate(0, 0, -2).Unix()
	to := time.Now().Unix()

	mockService.On("MMSBuscar", "BRLBTC", 20, from, to).Return([]dto.MMSResponse{}, errors.New("erro interno"))

	req, _ := http.NewRequest("GET", "/BRLBTC/mms?range=20&from="+
		strconv.FormatInt(from, 10)+"&to="+strconv.FormatInt(to, 10), nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "erro interno")
}
