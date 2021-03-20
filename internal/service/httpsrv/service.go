package httpsrv

import (
	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/middlewares"
	"github.com/gophers0/storage/internal/repository/postgres"
	"github.com/gophers0/storage/internal/service/httpsrv/handlers"
	"github.com/gophers0/storage/pkg/bindings"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	gaarx "github.com/zergu1ar/Gaarx"
	"net/http"
)

type Service struct {
	log  *logrus.Logger
	app  *gaarx.App
	name string
}

func New(log *logrus.Logger) *Service {
	return &Service{
		log:  log,
		name: "HTTPService",
	}
}

func (s *Service) GetName() string {
	return s.name
}

func (s *Service) getLog() *logrus.Entry {
	return s.log.WithField("service", s.name)
}

func (s *Service) Start(a *gaarx.App) error {
	s.app = a
	e := echo.New()
	e.Validator = &bindings.Validator{}

	mw := middlewares.New(a.Config(), a.GetLog(), a.GetDB().(*postgres.Repo))
	commonMw := []echo.MiddlewareFunc{
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowMethods: []string{"*"},
		}),
		mw.Log(),
		mw.Error(middlewares.ErrorHandler()),
		mw.Recover(),
	}
	authMw := append(commonMw, mw.Auth())
	//adminMw := append(authMw, mw.AdminOnly())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}, commonMw...)

	h := handlers.New(a)
	profile := e.Group("/profile", authMw...)
	{
		profile.GET("/", h.GetProfile)
		profile.OPTIONS("/", echo.MethodNotAllowedHandler)
	}

	return e.Start(":" + s.app.Config().(*config.Config).Api.Port)
}

func (s *Service) Stop() {
	s.log.Debug("stop" + s.name)
}
