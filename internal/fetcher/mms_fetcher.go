package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Junkes887/MMS/internal/constants"
	"github.com/Junkes887/MMS/internal/database/entity"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type CandleResponse struct {
	Close      []string `json:"c"`
	Timestamps []int64  `json:"t"`
}

var fetchCandles = func(symbol string, from, to int64) ([]float64, []int64, error) {
	url := fmt.Sprintf("https://api.mercadobitcoin.net/api/v4/candles?symbol=%s&from=%d&to=%d&resolution=1d", symbol, from, to)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept-Encoding", "identity")

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var data CandleResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, nil, err
	}

	var prices []float64

	for i := 0; i < len(data.Close) && i < len(data.Timestamps); i++ {
		price, err := strconv.ParseFloat(data.Close[i], 64)
		if err == nil {
			prices = append(prices, price)
		}
	}

	return prices, data.Timestamps, nil
}

func calcMMS(prices []float64, window int) []float64 {
	var result []float64
	for i := range prices {
		if i+1 < window {
			result = append(result, 0)
			continue
		}
		sum := 0.0
		for j := i - window + 1; j <= i; j++ {
			sum += prices[j]
		}
		result = append(result, sum/float64(window))
	}
	return result
}

func getDelay() time.Duration {
	delayStr := os.Getenv("DELAY_FETCHER_SECONDS")
	seconds, err := strconv.Atoi(delayStr)

	if err != nil {
		seconds = 120
	} else if seconds < 0 {
		seconds = 120
	}

	return time.Duration(seconds) * time.Second
}

func (f *Fetcher) SeedData(symbol string, pair string, from int64, to int64) {
	var prices []float64
	var timestamps []int64
	var err error

	maxRetries := 5
	delay := getDelay()

	for attempt := 1; attempt <= maxRetries; attempt++ {
		prices, timestamps, err = fetchCandles(symbol, from, to)
		if err == nil {
			break
		}

		f.Logger.Warn("Erro ao buscar candles, tentativa de novo em breve...",
			zap.Int("tentativa", attempt),
			zap.Error(err),
		)
		time.Sleep(delay)
	}

	if err != nil {
		f.Logger.Error("Erro ao buscar candles: ", zap.Error(err))
		return
	}

	mlen := len(prices)
	mms20 := calcMMS(prices, 20)
	mms50 := calcMMS(prices, 50)
	mms200 := calcMMS(prices, 200)

	for i := 0; i < mlen; i++ {
		f.Repository.SaveMSS(entity.MMSEntity{
			Pair:      pair,
			Timestamp: timestamps[i],
			MMS20:     mms20[i],
			MMS50:     mms50[i],
			MMS200:    mms200[i],
		})
	}

	f.Logger.Info("Dados inseridos com sucesso",
		zap.String("pair", pair),
		zap.Int("quantidade", mlen),
	)
}

func (f *Fetcher) RunDailyJob() {
	go func() {
		for {
			f.Logger.Info("Iniciando job de incremento diário")

			from := time.Now().AddDate(0, 0, -1).Unix()
			to := time.Now().Unix()
			for symbol, pair := range constants.SymbolPairMap {
				f.SeedData(symbol, pair, from, to)
			}

			f.Logger.Info("Job finalizado, aguardando 24h...")
			time.Sleep(24 * time.Hour)
		}
	}()
}

func (f *Fetcher) VerificarDadosFaltantes() {
	now := time.Now().Truncate(24 * time.Hour)
	start := now.AddDate(0, 0, -365).Unix()
	end := now.Unix()

	for _, pair := range constants.SymbolPairMap {
		existentes := f.Repository.BuscarDiasFaltantes(pair, start, end)

		existsMap := make(map[int64]bool)
		for _, ts := range existentes {
			existsMap[ts] = true
		}

		var faltando []int64
		for ts := start; ts <= end; ts += 86400 {
			if !existsMap[ts] {
				faltando = append(faltando, ts)
			}
		}

		if len(faltando) > 0 {
			f.Logger.Warn("Dias faltando detectados", zap.String("pair", pair), zap.Int("quantidade", len(faltando)))
			enviarEmailAlerta(f, pair, faltando)
		} else {
			f.Logger.Info("Nenhum dia faltando para o par",
				zap.String("pair", pair),
				zap.Int("dias_completos", len(existentes)),
			)
		}
	}
}

var enviarEmailAlerta = func(f *Fetcher, pair string, diasFaltando []int64) {
	m := gomail.NewMessage()
	m.SetHeader("From", "vitrine.me.site@gmail.com")
	m.SetHeader("To", os.Getenv("EMAIL_ALERT"))
	m.SetHeader("Subject", fmt.Sprintf("ALERTA: Falha em dados de %s", pair))
	body := fmt.Sprintf("Foram encontrados %d dias faltando para o par %s:\n\n", len(diasFaltando), pair)
	for _, ts := range diasFaltando {
		dia := time.Unix(ts, 0).Format("02/01/2006")
		body += fmt.Sprintf("- %s\n", dia)
	}
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "vitrine.me.site@gmail.com", os.Getenv("EMAIL_SENHA"))

	if err := d.DialAndSend(m); err != nil {
		f.Logger.Error("Erro ao enviar e-mail de alerta", zap.Error(err))
	}
}
