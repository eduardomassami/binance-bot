package main

import (
	"context"
	"log"
	"os"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/eduardomassami/binance-bot/config"
	"github.com/joho/godotenv"
)

func main() {
	logger := config.GetLogger("main")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
  
  client := binance_connector.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))
  client.BaseURL = "https://testnet.binance.vision/api"

  err = client.NewPingService().Do(context.Background())
	if err != nil {
		log.Fatalf("Erro ao conectar na Binance: %v", err)
	}

	logger.Debug("Conexão bem-sucedida")

  getMarketPrice(client, "BTCUSDT")
  simpleTrade(client, "BTCUSDT", 30000.0, 35000.0)
}

func getMarketPrice(client *binance_connector.Client, symbol string) {
	logger := config.GetLogger("main")
	price, err := client.NewAvgPriceService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Fatalf("Erro ao obter preço do mercado: %v", err)
	}
	logger.Debugf("Preço médio de %s: %s", symbol, price.Price)
}

func buyOrder(client *binance_connector.Client, symbol string, quantity float64) {
	logger := config.GetLogger("main")
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side("Buy").
		Type("Market").
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Erro ao realizar compra: %v", err)
	}
  logger.Debugf("Ordem de compra realizada: %v\n", order)
}

func sellOrder(client *binance_connector.Client, symbol string, quantity float64) {
	logger := config.GetLogger("main")
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side("Sell").
		Type("Market").
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Erro ao realizar venda: %v", err)
	}
	logger.Debugf("Ordem de venda realizada: %v\n", order)
}

func simpleTrade(client *binance_connector.Client, symbol string, buyThreshold float64, sellThreshold float64) {
	logger := config.GetLogger("main")
	price, _ := client.NewAvgPriceService().Symbol(symbol).Do(context.Background())

	currentPrice, _ := strconv.ParseFloat(price.Price, 64)

	if currentPrice < buyThreshold {
		buyOrder(client, symbol, 0.001)
	} else if currentPrice > sellThreshold {
		sellOrder(client, symbol, 0.001)
	} else {
		logger.Debug("Nenhuma ação necessária.")
	}
}
