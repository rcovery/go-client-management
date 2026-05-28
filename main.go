package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rcovery/go-client-management/internal/config"
	"github.com/rcovery/go-client-management/internal/http/handlers"
	"github.com/rcovery/go-client-management/internal/infra/postgres"
)

func main() {
	config.InitConfig()

	baseCtx := context.Background()
	connectionString := postgres.GetConnectionFromEnv()
	db, databaseErr := postgres.NewDatabaseConnection(connectionString)
	if databaseErr != nil {
		panic(databaseErr)
	}

	host := config.GetString("HOST")
	port := config.GetString("PORT")

	log.Printf("Listen %s:%s", host, port)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		BaseContext:  func(net.Listener) context.Context { return baseCtx },
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	}

	h := handlers.New(db)
	h.HandleClient()
	// h.HandleWebhook()

	server.ListenAndServe()
}
