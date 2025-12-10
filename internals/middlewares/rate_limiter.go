package middlewares

import "net/http"

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do rate limit
		next.ServeHTTP(w, r)
	})
}
