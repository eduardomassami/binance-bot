package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	logger.Debugf("apiKey: %s and secretKey: %s", apiKey, secretKey)

	// ----------------------------------------------------------------------------------------------------------------------

	// Initialise Websocket API Client
	client := binance_connector.NewWebsocketAPIClient(apiKey, secretKey)
	// Connect to Websocket API
	err = client.Connect()
	if err != nil {
		logger.Errorf("Error connecting: %v", err)
		return
	}
	defer client.Close()

	// Send request to Websocket API
	response, err := client.NewAccountOCOHistoryService().Do(context.Background())
	if err != nil {
		logger.Errorf("Error send request: %v", err)
		return
	}

	// Print the response
	fmt.Println(binance_connector.PrettyPrint(response))

	client.WaitForCloseSignal()
}
