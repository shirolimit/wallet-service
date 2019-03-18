package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shirolimit/wallet-service/pkg/db"
	"github.com/shirolimit/wallet-service/pkg/service"

	log "github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/shirolimit/wallet-service/pkg/endpoint"
	"github.com/shirolimit/wallet-service/pkg/transport"
)

var (
	fs       = flag.NewFlagSet("wallet", flag.ExitOnError)
	httpAddr = fs.String("http-address", ":8080", "HTTP address to listen")
	connStr  = fs.String("connection-string", "", "Postgres connection string")
	// zipkinURL = fs.String("zipkin-url", "", "URL for Zipkin tracing")
)

func main() {
	fs.Parse(os.Args[1:])

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)

	storage := db.NewPgStorage(*connStr)

	svc := service.NewWalletService(storage)
	svc = service.LoggingMiddleware(logger)(svc)

	endpoints := endpoint.NewEndpointSet(svc)

	handler := transport.NewHTTPHandler(endpoints, nil)
	server := http.Server{Addr: *httpAddr, Handler: handler}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log(
				"transport", "HTTP",
				"error", err,
			)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
