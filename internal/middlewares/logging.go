package middlewares

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Logging(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.WithFields(logrus.Fields{
				"remote_addr": r.RemoteAddr,
			})
			logger.Infof("started %s %s", r.Method, r.RequestURI)

			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Infof("completed in %v", time.Since(start))
		})
	}
}
