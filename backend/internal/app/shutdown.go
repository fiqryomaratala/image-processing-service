package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/database"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/queue"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/server"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/storage"
	"go.uber.org/zap"
)

const shutdownTimeout = 10 * time.Second

func Run(httpServer *server.Server) error {
	log := logger.Get()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	var serverErrCh <-chan error
	if httpServer != nil {
		errCh := make(chan error, 1)
		serverErrCh = errCh

		go func() {
			if err := httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
			close(errCh)
		}()
	}

	select {
	case sig := <-signalCh:
		log.Info("Shutdown signal received", zap.String("signal", sig.String()))
		return shutdown(httpServer)
	case err, ok := <-serverErrCh:
		if !ok {
			return nil
		}

		return err
	}
}

func shutdown(httpServer *server.Server) error {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var shutdownErr error

	if httpServer != nil {
		log.Info("Stopping HTTP server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error("Failed to stop HTTP server", zap.Error(err))
			shutdownErr = errors.Join(shutdownErr, err)
		} else {
			log.Info("HTTP server stopped")
		}
	}

	log.Info("Closing PostgreSQL connection...")
	if err := database.Close(); err != nil {
		log.Error("Failed to close PostgreSQL connection", zap.Error(err))
		shutdownErr = errors.Join(shutdownErr, err)
	} else {
		log.Info("PostgreSQL connection closed")
	}

	log.Info("Closing RabbitMQ channel...")
	if err := queue.CloseChannel(); err != nil {
		log.Error("Failed to close RabbitMQ channel", zap.Error(err))
		shutdownErr = errors.Join(shutdownErr, err)
	} else {
		log.Info("RabbitMQ channel closed")
	}

	log.Info("Closing RabbitMQ connection...")
	if err := queue.CloseConnection(); err != nil {
		log.Error("Failed to close RabbitMQ connection", zap.Error(err))
		shutdownErr = errors.Join(shutdownErr, err)
	} else {
		log.Info("RabbitMQ connection closed")
	}

	log.Info("Releasing MinIO resources...")
	if err := storage.Close(); err != nil {
		log.Error("Failed to release MinIO resources", zap.Error(err))
		shutdownErr = errors.Join(shutdownErr, err)
	} else {
		log.Info("MinIO resources released")
	}

	if shutdownErr != nil {
		log.Warn("Shutdown completed with errors", zap.Error(shutdownErr))
		return shutdownErr
	}

	log.Info("Shutdown completed successfully")

	return nil
}
