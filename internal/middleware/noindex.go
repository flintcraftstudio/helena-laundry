package middleware

import "net/http"

// NoIndex returns middleware that stamps an X-Robots-Tag: noindex, nofollow
// header on every response when enabled — the most robust crawl signal, since it
// covers non-HTML responses too. A no-op when disabled, so production is
// unaffected. Pair with a Disallow-all robots.txt for testing/staging deploys.
func NoIndex(enabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !enabled {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Robots-Tag", "noindex, nofollow")
			next.ServeHTTP(w, r)
		})
	}
}
