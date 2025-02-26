package api

import (
	"applicationDesignTest/internal/api/handlers"
	"applicationDesignTest/internal/services"
	"applicationDesignTest/pkg/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	orderHandler   *handlers.OrderHandler
	bookingManager services.BookingManager
	logger         logger.Logger
}

func NewServer(bookingManager services.BookingManager, logger logger.Logger) *Server {
	router := gin.Default()

	server := &Server{
		router:         router,
		bookingManager: bookingManager,
		orderHandler:   handlers.NewOrderHandler(bookingManager),
		logger:         logger,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	s.orderHandler.RegisterRoutes(s.router)
}

func (s *Server) Run(port string) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server: %v", err)
		}
	}()

	s.logger.Info("Server is running on port %s", port)

	<-quit
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server forced to shutdown: %v", err)
	}

	s.logger.Info("Server exited properly")
	return nil
}
