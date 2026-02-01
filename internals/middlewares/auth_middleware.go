package middlewares

import (
	"context"
	"file-analyzer/internals/services"
	"file-analyzer/internals/utils"
	"net/http"
	"strconv"
	"strings"
)

type UserID struct{}

func Auth(s services.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.FAIL(w, http.StatusUnauthorized, "Missing Authorization Header")
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.FAIL(w, http.StatusUnauthorized, "Invalid Authorization Header")
				return
			}
			userID, err := s.VerifyToken(parts[1])
			if err != nil {
				utils.FAIL(w, http.StatusUnauthorized, "Invalid token")
				return
			}
			userIDInt,err := strconv.ParseInt(userID,10,64)
			if err != nil {
				utils.FAIL(w, http.StatusUnauthorized, "Error converting User ID")
				return
			}
			ctx := context.WithValue(r.Context(), UserID{}, userIDInt)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
