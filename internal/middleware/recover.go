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
				case string, error:
					err = fmt.Errorf("panic: %s", errType)
				default:
					err = errors.New("unknown error")
				}
				m.l.Error(err, stackTrace)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
