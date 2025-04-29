package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) MMSBuscar(c *gin.Context) {
	pair := c.Param("pair")
	rangeParam := c.Query("range")
	fromStr := c.Query("from")
	toStr := c.DefaultQuery("to", strconv.FormatInt(time.Now().AddDate(0, 0, -1).Unix(), 10))

	mmsRange, err := strconv.Atoi(rangeParam)
	if err != nil || (mmsRange != 20 && mmsRange != 50 && mmsRange != 200) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "range inválido, deve ser 20, 50 ou 200"})
		return
	}

	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from inválido"})
		return
	}

	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to inválido"})
		return
	}

	resp, err := h.Service.MMSBuscar(pair, mmsRange, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Info("Busca de MMS realizada com sucesso",
		zap.String("pair", pair),
		zap.Int("range", mmsRange),
		zap.Int("quantidade_resultados", len(resp)),
	)

	c.JSON(http.StatusOK, resp)
}
