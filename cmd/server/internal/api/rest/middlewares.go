package rest

import (
	"net"
	"net/http"

	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
	"github.com/stas-zatushevskii/Monitor/cmd/server/internal/logger"
)

func WhiteListMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	if cfg.TrustedSubnet == "" {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, trustedNet, err := net.ParseCIDR(cfg.TrustedSubnet)
			if err != nil {
				logger.Log.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			host := r.Header.Get("X-Real-IP")

			ip := net.ParseIP(host)
			if ip == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if !trustedNet.Contains(ip) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
