package middleware

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	written *bool
}

func (rw *responseWriter) WriteHeader(code int) {
	*rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	*rw.written = true
	return rw.ResponseWriter.Write(b)
}

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var written bool
		mw := &responseWriter{w, &written}

		defer func() {
			if err := recover(); err != nil {
				if !written {
					http.Error(mw, "Internal Server Error", http.StatusInternalServerError)
				}
				log.Printf("Error: %v", err)
			}
		}()

		next.ServeHTTP(mw, r)
	})
}
