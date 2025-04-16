package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type ctxKeyRequestID struct{}

type responseWriter struct {
	http.ResponseWriter
	written *bool
	status  int
}

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func generateRequestId() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	// Превращаем байты в строку
	randomString := hex.EncodeToString(bytes)
	return randomString
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
		requestID := generateRequestId()
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, requestID)
		r = r.WithContext(ctx)

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
				"request_id": requestID,
			}).Info("Request")
		} else {
			log.Println("Logger is nil")
		}
	})
}
