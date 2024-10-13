package http

import (
	"net/http"
)

func APIKeyAuthMiddleware(APIKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("Authorization")

			if apiKey == "" || apiKey != APIKey {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
