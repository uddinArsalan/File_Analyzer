package middlewares

import (
	"file-analyzer/internals/utils"
	"log"
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Println()
			utils.FAIL(w, http.StatusInternalServerError, "Failed to do Basic Auth ")
			return
		}
		if username != "Arsu" || password != "pass" {
			utils.FAIL(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next(w, r)
	}
}
