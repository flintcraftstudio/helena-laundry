package handler

import "net/http"

// Robots serves /robots.txt. When noindex is set (a testing/staging deploy) it
// disallows all crawling; otherwise it allows everything.
func Robots(noindex bool) http.HandlerFunc {
	body := "User-agent: *\nDisallow:\n"
	if noindex {
		body = "User-agent: *\nDisallow: /\n"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}
}
