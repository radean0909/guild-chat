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

	"github.com/radean0909/guild-chat/api/handlers"
	"github.com/radean0909/guild-chat/api/internal/db"
	"github.com/radean0909/guild-chat/api/internal/db/mem"
)

// HealthGracePeriod provided enough time to wait after a SIGTERM has been
// provided to verify a healthcheck has been made against the service to update
// kubernetes load balancers so no traffic is routed to the service during
// graceful shutdown.
const HealthGracePeriod time.Duration = 10 * time.Second

type Service struct {
	echo         *echo.Echo
	DB           db.Driver
	MsgHandler   *handlers.MessageHandler
	ConvoHandler *handlers.ConversationHandler
	UserHandler  *handlers.UserHandler
	ready        bool
}

func New() *Service {
	s := &Service{}
	e := echo.New()
	e.HideBanner = false
	e.HidePort = false

	// Set up the database, for the example we are using an in memory datastore
	// In production, this would connect to something more permanent
	// Depending on more specific details (more reads than writes, future features)
	// I would likely choosd Postgres or MongoDB
	s.DB = mem.NewDriver()

	// Set up the individual "handlers" - in production this might dial out to gRPC handler services
	s.MsgHandler = &handlers.MessageHandler{
		DB: s.DB,
	}

	s.ConvoHandler = &handlers.ConversationHandler{
		DB: s.DB,
	}

	s.UserHandler = &handlers.UserHandler{
		DB: s.DB,
	}

	// logger - in production this would likely be more robust
	e.Use(middleware.Logger())

	// global middleware
	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	// authentication/authorization middlewares could exist at the top level or on individual groups or routes

	// kubernetes health checks - important for pod green status when deployed in a container on the cloud
	e.GET("/alive", s.HandleAlive())
	e.GET("/ready", s.HandleReady())

	// system status - useful for local api testing, otherwise, could be used to provide other metadata
	// such as messages sent, uptime, etc
	e.GET("/system/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "It's alive!")
	})

	// message endpoints - singular message between two users
	msgs := e.Group("/message")
	msgs.POST("", s.postMessage)
	msgs.GET("/:id", s.getMessageByID)

	// converstion endpoints - a conversation includes all messages between two users
	conversations := e.Group("/conversation")
	conversations.GET("/:to/:from", s.getConversation)
	conversations.GET("/:to", s.listConversations)

	// user endpoints
	users := e.Group("/user")
	users.POST("", s.postUser)
	users.GET("/:id", s.getUserByID)
	users.DELETE("/:id", s.deleteUserByID)

	s.echo = e

	return s
}

// Start the Service listening on addr. On a SIGTERM the Service will
// start a graceful shutdown.
func (s *Service) Start(addr string) error {
	s.ready = true
	defer func() { s.ready = false }()

	// gracefully stop on SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		s.echo.Logger.Info("shutting down...")
		s.ready = false
		time.Sleep(HealthGracePeriod)
		s.echo.Shutdown(context.Background())
	}()

	return s.echo.Start(addr)
}

// HandleReady checks for kubernetes
func (s *Service) HandleReady() echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.ready {
			return c.NoContent(http.StatusOK)
		}
		return c.NoContent(http.StatusServiceUnavailable)
	}
}

// HandleAlive checks for kubernetes
func (s *Service) HandleAlive() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}

// root handlers - these are used to help pass DB down to the lower level handler
// messages
func (s *Service) getMessageByID(c echo.Context) error {
	return s.MsgHandler.GetMessageByID(c)
}

func (s *Service) postMessage(c echo.Context) error {
	return s.MsgHandler.PostMessage(c)
}

// conversations
func (s *Service) getConversation(c echo.Context) error {
	return s.ConvoHandler.GetConversation(c)
}

func (s *Service) listConversations(c echo.Context) error {
	return s.ConvoHandler.ListConversations(c)
}

// users
func (s *Service) postUser(c echo.Context) error {
	return s.UserHandler.PostUser(c)
}

func (s *Service) getUserByID(c echo.Context) error {
	return s.UserHandler.GetUserByID(c)
}

func (s *Service) deleteUserByID(c echo.Context) error {
	return s.UserHandler.DeleteUserbyID(c)
}
