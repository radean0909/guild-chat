package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"

	"github.com/radean0909/guild-chat/handlers"
	"github.com/radean0909/guild-chat/internal/db"
)


// HealthGracePeriod provided enough time to wait after a SIGTERM has been
// provided to verify a healthcheck has been made against the service to update
// kubernetes load balancers so no traffic is routed to the service during
// graceful shutdown.
const HealthGracePeriod time.Duration = 10 * time.Second

type Server struct {
	echo  *echo.Echo
	db *db.Driver
	ready bool
}

func New(log *logrus.Entry) *Server {
	s := &Server{}
	e := echo.New()
	e.HideBanner = false
	e.HidePort = false

	e.Logger = echologrus.Logger{Logger: log.Logger}
	e.Use(logging.RequestLogs(log))

	// global middleware
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	// kubernetes health checks - important for pod green status when deployed in a container on the cloud
	e.GET("/alive", s.HandleAlive())
	e.GET("/ready", s.HandleReady())

	// system status - useful for local api testing, otherwise, could be used to provide other metadata
	// such as messages sent, uptime, etc
	e.GET("/system/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// message endpoints - singular message between two users
	msgs := e.Group("/message")
	msgs.GET("/:id", handlers.GetMessageByID)
	msgs.POST("/", handlers.PostMessage)

	// converstion endpoings - a conversation includes all messages between two users
	conversations := e.Group("/conversation")
	conversations.GET("/:to/:from", handlers.GetConversation)
	conversations.GET("/:to", GetAllConversations)

	users := e.Group("/user")
	users.POST("/", handlers.PostUser)
	users.GET("/:id", handlers.GetUserByID)
	users.DELETE("/:id", handlers.DeleteUserByID)

	e.GET("/check-in/:extSysId", handlers.GetCheckIn)
	e.POST("/check-in", handlers.PostCheckIn)
	s.echo = e

	return s
}

// Start the server listening on addr. On a SIGTERM the server will
// start a graceful shutdown.
func (s *Server) Start(addr string) error {
	s.ready = true
	defer func() { s.ready = false }()

	// gracefully stop on SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		s.ready = false
		time.Sleep(HealthGracePeriod)
		s.echo.Shutdown(context.Background())
	}()

	return s.echo.Start(addr)
}

// HandleReady checks for kubernetes
func (s *Server) HandleReady() echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.ready {
			return c.NoContent(http.StatusOK)
		}
		return c.NoContent(http.StatusServiceUnavailable)
	}
}

// HandleAlive checks for kubernetes
func (s *Server) HandleAlive() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}