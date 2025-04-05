package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
}

type responseWriter struct {
	http.ResponseWriter
	written *bool
	status  int
}

func (rw *responseWriter) WriteHeader(code int) {
	*rw.written = true
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	*rw.written = true
	return rw.ResponseWriter.Write(b)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var written bool
		mw := &responseWriter{w, &written, 0}

		defer func() {
			if err := recover(); err != nil {
				if !written {
					http.Error(mw, "Internal Server Error", http.StatusInternalServerError)
				}
				if logger != nil {
					logger.WithFields(logrus.Fields{
						"error": err,
					}).Error("Recovered from panic")
				} else {
					log.Println("Logger is nil")
				}
			}
		}()

		startTime := time.Now()
		next.ServeHTTP(mw, r)
		latency := time.Since(startTime)

		if logger != nil {
			logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"url":        r.URL,
				"status":     mw.status,
				"latency":    latency,
				"user_agent": r.UserAgent(),
			}).Info("Request")
		} else {
			log.Println("Logger is nil")
		}
	})
}
