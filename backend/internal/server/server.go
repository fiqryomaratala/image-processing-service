package server

import (
	"context"
	"net/http"
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
}

func New() *Server {
	cfg := config.Get()

	if cfg.App.Env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.RequestID(),
		middleware.CORS(),
		middleware.Logger(),
		middleware.Recovery(),
	)

	registerRoutes(engine)

	return &Server{
		engine: engine,
		http: &http.Server{
			Addr:              cfg.App.Address(),
			Handler:           engine,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	log := logger.Get()
	log.Info("Listening on "+s.http.Addr, zap.String("address", s.http.Addr))

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
