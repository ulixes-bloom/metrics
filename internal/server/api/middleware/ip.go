package middleware

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
)

func WithIPResolving(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ipReq, err := resolveIP(r)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, ipnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				log.Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !ipnet.Contains(ipReq) {
				log.Info().Msgf("Untrusted IP address: %s", ipReq)
				http.Error(w, "Untrusted IP address", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func resolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get(headers.XRealIP)
	ip := net.ParseIP(ipStr)

	if ip == nil {
		return nil, fmt.Errorf("ip.resolveIP: failed to parse ip from http header")
	}
	return ip, nil
}
