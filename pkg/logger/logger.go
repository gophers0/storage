package logger

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// Keys used in logger.
const (
	InstanceKey     = "instance-key"
	KeyRequest      = "request"
	KeyResponse     = "response"
	KeyRequestID    = "request-id"
	KeyRequestDump  = "request-dump"
	KeyResponseDump = "response-dump"
	KeyError        = "error"
	KeyPanic        = "panic"
	KeySecondPanic  = "second-panic" // panic throw while logging first panic
)

// Keys used in context.
const (
	ContextLoggerKey = "main.logger"
)

func AddField(ctx echo.Context, fildeName, value string) {
	log := ctx.Get(ContextLoggerKey).(*logrus.Entry)
	ctx.Set(ContextLoggerKey, log.WithField(fildeName, value))
}

// logger returns logrus Entry stored in context, if any.
func Logger(ctx echo.Context) *logrus.Entry {
	if log, ok := ctx.Get(ContextLoggerKey).(*logrus.Entry); ok {
		return log
	}

	return nil
}

// LogFromContext returns logrus Entry stored in context, if any.
func LogFromContext(ctx echo.Context) (*logrus.Entry, bool) {
	log, ok := ctx.Get(ContextLoggerKey).(*logrus.Entry)
	return log, ok
}

// filterHeader returns filtered copy of header.
func FilterHeader(header http.Header) http.Header {
	newHeader := make(http.Header)
	for k, v := range header {
		switch k {
		case "Authorization":
		// ignore these headers
		default:
			newHeader[k] = v
		}
	}

	return newHeader
}

// httpRequest is used to log http transport.
type HttpRequest struct {
	Method     string      `json:"method"`
	Header     http.Header `json:"-"`
	RequestURI string      `json:"request-uri"`
	RemoteAddr string      `json:"remote-addr"`
	URL        *url.URL    `json:"-"`
}

// httpResponse is used to log http responses.
type HttpResponse struct {
	Status           int     `json:"status"`
	Size             int64   `json:"size"`
	Duration         string  `json:"duration"`
	WriteDuration    string  `json:"write-duration"`
	DurationSec      float64 `json:"duration-sec"`
	WriteDurationSec float64 `json:"write-duration-sec"`
}
