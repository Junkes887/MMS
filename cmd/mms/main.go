package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Junkes887/MMS/internal/api"
	"github.com/Junkes887/MMS/internal/constants"
	"github.com/Junkes887/MMS/internal/database"
	"github.com/Junkes887/MMS/internal/database/repository"
	"github.com/Junkes887/MMS/internal/fetcher"
	"github.com/Junkes887/MMS/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	port := os.Getenv("PORT")
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	conn := database.NewConfig()
	repository := repository.NewRepository(conn, logger)
	service := service.NewService(repository, logger)
	handler := api.NewHandler(service, logger)
	fetcher := fetcher.NewFetcher(repository, logger)

	router := gin.Default()
	router.GET("/:pair/mms", handler.MMSBuscar)

	from := time.Now().AddDate(0, 0, -365).Unix()
	to := time.Now().Unix()

	for symbol, pair := range constants.SymbolPairMap {
		go fetcher.SeedData(symbol, pair, from, to)
	}
	go fetcher.RunDailyJob()
	go fetcher.VerificarDadosFaltantes()

	router.Run(fmt.Sprintf(":%s", port))
}
