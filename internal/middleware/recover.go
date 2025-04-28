package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (m *Middleware) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var err error
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				switch errType := r.(type) {
				case string:
					err = errors.New(fmt.Sprintf("panic: %s", errType))
				case error:
					err = errors.New(fmt.Sprintf("panic: %s", errType))
				default:
					err = errors.New("Unknown error")
				}
				m.l.Error("panic:", err, stackTrace)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
