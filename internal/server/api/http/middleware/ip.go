package middleware

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/headers"
)

// WithIPResolving is a middleware that checks whether the client's IP address is within a trusted subnet.
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

// resolveIP extracts and parses the IP address from the "X-Real-IP" header in the HTTP request.
// If the header is missing or the IP is invalid, it returns an error.
func resolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get(headers.XRealIP)
	ip := net.ParseIP(ipStr)

	if ip == nil {
		return nil, fmt.Errorf("ip.resolveIP: failed to parse ip from http header")
	}
	return ip, nil
}
