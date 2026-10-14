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
	"github.com/nashinnov8/go-ecommerce/services/product/internal/config"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/repository/postgres"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/service"
	httptransport "github.com/nashinnov8/go-ecommerce/services/product/internal/transport/http"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/transport/http/handler"
)

func main() {
	// 1. Load Configuration & Initialize Structured Logger
	cfg := config.Load()
	log := logger.New(cfg.Environment)

	log.Info("starting service",
		slog.String("service", cfg.ServiceName),
		slog.String("env", cfg.Environment),
		slog.String("port", cfg.HTTPPort),
	)

	// 2. Initialize PostgreSQL Connection Pool (pgx/v5)
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
	log.Info("connected to postgres database")

	// 3. Wire Dependencies (Repository -> Service -> Handler -> Router)
	productRepo := postgres.NewProductRepository(dbPool)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)
	router := httptransport.NewRouter(log, cfg.Environment, productHandler)

	// 4. Configure HTTP Server with Timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start HTTP Server in a Background Goroutine
	go func() {
		log.Info("http server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// 6. Wait for Interrupt Signal for Graceful Shutdown
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
