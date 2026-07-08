package server

import (
	"context"
	"net/http"
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/handler"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
	log    *zap.Logger
}

func New(appCfg config.AppConfig, corsCfg config.CORSConfig, log *zap.Logger, healthHandler *handler.HealthHandler) *Server {
	if appCfg.Env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.RequestID(),
		middleware.CORS(corsCfg),
		middleware.Logger(log),
		middleware.Recovery(log),
	)

	registerRoutes(engine, healthHandler)

	return &Server{
		engine: engine,
		http: &http.Server{
			Addr:              appCfg.Address(),
			Handler:           engine,
			ReadHeaderTimeout: 5 * time.Second,
		},
		log: log,
	}
}

func (s *Server) Run() error {
	s.log.Info("Listening on "+s.http.Addr, zap.String("address", s.http.Addr))

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
