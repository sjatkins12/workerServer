package utils

/* This is the logging module in the utils package. We setup logging based on the type of Environment we are dealing with
 */
import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//Log - The principan logging interface that we are going to export
var (
	Log *logrus.Entry
)

// GetLogger turns a logrus logger with appropriate log level depending on environment
func GetLogger(isProd bool) *logrus.Logger {
	logger := logrus.New()
	if isProd {
		logger.Level = logrus.DebugLevel
		logger.Formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
			TimestampFormat: time.RFC3339Nano,
		}
	} else {
		logger.Level = logrus.DebugLevel
		logger.Formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339Nano}
	}
	return logger
}

//init is auto called by go runtime to setup things. We use it to setup logging here.
func init() {
	logger := GetLogger(true)
	Log = logger.WithFields(logrus.Fields{
		"env": "testing",
	})
}

// Source:	https://github.com/gin-gonic/contrib/blob/master/ginrus/ginrus.go

type loggerEntryWithFields interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

// Ginrus returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func Ginrus(logger loggerEntryWithFields) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path

		buf, _ := ioutil.ReadAll(c.Request.Body)

		var body interface{}
		err := json.Unmarshal(buf, &body)

		if err != nil {
			body = string(buf)
		}

		reader := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = reader

		c.Next()

		timeFormat := "2006-01-02T15:04:05.999999"

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		method := c.Request.Method
		status := c.Writer.Status()
		message := method + " " + path + " " + strconv.Itoa(status)

		entry := logger.WithFields(logrus.Fields{
			"status":     status,
			"method":     method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(timeFormat),
			"body":       body,
		})

		sessionID := c.Params.ByName("session_id")
		if sessionID != "" {
			entry.WithField("session_id", sessionID)
		}

		elementID := c.Params.ByName("element_id")
		if elementID != "" {
			entry.WithField("element_id", elementID)
		}

		if status >= 500 {
			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				entry = entry.WithField("errors", c.Errors.String())
			}
			entry.Error(message)
		} else if status >= 400 && status < 500 {
			entry.Warning(message)
		} else {
			entry.Info(message)
		}
	}

}
