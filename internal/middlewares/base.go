package middlewares

import (
	"github.com/gophers0/storage/pkg/users"
	"time"

	"github.com/google/uuid"
	"github.com/gophers0/storage/internal/config"
	"github.com/gophers0/storage/internal/repository/postgres"
	"github.com/gophers0/storage/pkg/logger"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	repo    *postgres.Repo
	cfg     *config.Config
	logger  *logrus.Logger
	UserApi *users.Api
}

func New(cfg interface{}, logger *logrus.Logger, repo *postgres.Repo) *Middleware {
	appCfg := cfg.(*config.Config)
	return &Middleware{
		repo:    repo,
		cfg:     appCfg,
		logger:  logger,
		UserApi: users.NewApi(appCfg.Users),
	}
}

// Log is a middleware which adds logger to context and logs every request.
func (mw *Middleware) Log() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			requestID := uuid.New().String()
			ctx.Set(logger.KeyRequestID, requestID)

			fields := logrus.Fields{
				logger.InstanceKey:  mw.cfg.Name,
				logger.KeyRequestID: requestID,
				logger.KeyRequest: logger.HttpRequest{
					Method:     ctx.Request().Method,
					Header:     logger.FilterHeader(ctx.Request().Header),
					RequestURI: ctx.Request().RequestURI,
					RemoteAddr: ctx.RealIP(),
					URL:        ctx.Request().URL,
				},
			}

			log := mw.logger.WithFields(fields)
			ctx.Set(logger.ContextLoggerKey, log)

			start := time.Now().UTC()

			var responseReadyTime time.Time
			ctx.Response().Before(func() {
				responseReadyTime = time.Now().UTC()
			})

			defer func() {
				res := ctx.Response()
				log = ctx.Get(logger.ContextLoggerKey).(*logrus.Entry)
				fields := logrus.Fields{
					logger.KeyResponse: logger.HttpResponse{
						Status:           res.Status,
						Size:             res.Size,
						Duration:         responseReadyTime.Sub(start).String(),
						DurationSec:      responseReadyTime.Sub(start).Seconds(),
						WriteDuration:    time.Since(responseReadyTime).String(),
						WriteDurationSec: time.Since(responseReadyTime).Seconds(),
					},
				}

				if ua := ctx.Request().Header.Get("User-Agent"); ua != "" {
					fields["user-agent"] = ua
				}
				log.WithFields(fields).Info("http response")
			}()

			log.WithFields(fields).Info("http request")
			return next(ctx)
		}
	}
}
