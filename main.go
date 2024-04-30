package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daniel-burghardt/ethereum-parser/data"
	"github.com/daniel-burghardt/ethereum-parser/ethrpc"
	"github.com/daniel-burghardt/ethereum-parser/httphandler"
	"github.com/daniel-burghardt/ethereum-parser/observer"
)

const serverPort = 3000

func main() {
	// Get ENVs
	serverUrl := os.Getenv("ETH_SERVER_URL")
	if serverUrl == "" {
		serverUrl = "https://cloudflare-eth.com"
	}
	webhookUrl := os.Getenv("WEBHOOK_URL")

	// Setup dependencies
	repo := data.NewInMemoryStorage()
	handler := httphandler.Handler{
		Repo: repo,
	}
	rpcService := ethrpc.Service{
		Url: serverUrl,
	}
	observerService := observer.Service{
		RPC:        rpcService,
		Repo:       repo,
		WebhookUrl: webhookUrl,
	}

	go observerService.Start()

	log.Printf("Starting the server on port %d...", serverPort)
	http.HandleFunc("/", http.NotFound)
	http.HandleFunc("GET /current-block", handler.GetCurrentBlock)
	http.HandleFunc("POST /subscribe/{address}", handler.PostSubscribe)
	http.HandleFunc("GET /transactions/{address}", handler.GetTransactions)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", serverPort), nil)
}
