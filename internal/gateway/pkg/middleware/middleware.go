package middleware

import (
	"log"
	"net/http"
	"time"
)

/* RequestLogger logs information about incoming HTTP requests. */
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		log.Printf("REQUEST: %s %s %s | Remote: %s | User-Agent: %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)

		next.ServeHTTP(w, r)

		log.Printf("RESPONSE: %s %s | Duration: %v",
			r.Method,
			r.URL.Path,
			time.Since(startTime),
		)
	})
}
