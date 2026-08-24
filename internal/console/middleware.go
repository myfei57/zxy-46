package console

import (
	"net/http"
	"time"

	"bldghvac/internal/audit"
)

func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		_ = s.deps.Audit.Record(audit.Entry{
			Source: "console",
			ZoneID: r.URL.Path,
			Action: r.Method,
			Detail: time.Since(start).String(),
		})
	})
}
