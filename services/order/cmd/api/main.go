package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nashinnov8/go-ecommerce/pkg/database"
	"github.com/nashinnov8/go-ecommerce/pkg/logger"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/config"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/repository/postgres"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/service"
	httptransport "github.com/nashinnov8/go-ecommerce/services/order/internal/transport/http"
	"github.com/nashinnov8/go-ecommerce/services/order/internal/transport/http/handler"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Environment)

	log.Info("starting service",
		slog.String("service", cfg.ServiceName),
		slog.String("port", cfg.HTTPPort),
	)

	// 1. Initialize PostgreSQL Pool
	ctx := context.Background()
	dbPool, err := database.NewPostgresPool(ctx, database.PostgresConfig{
		DSN:      cfg.PostgresDSN,
		MaxConns: cfg.DBMaxConns,
		MinConns: cfg.DBMinConns,
	})
	if err != nil {
		log.Error("failed to connect to postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	// 2. Wire Dependencies: TxManager -> Repo -> Service -> Handler -> Router
	txManager := database.NewTransactionManager(dbPool)
	orderRepo := postgres.NewOrderRepository(dbPool)
	orderSvc := service.NewOrderService(orderRepo, txManager)
	orderHandler := handler.NewOrderHandler(orderSvc)
	router := httptransport.NewRouter(log, cfg.Environment, orderHandler)

	// 3. HTTP Server with Graceful Shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("order-service listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutting down server...", slog.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server exited gracefully")
}
