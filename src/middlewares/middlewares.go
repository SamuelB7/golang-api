package middlewares

import (
	"api/src/auth"
	"api/src/responses"
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidateToken(r)
		if err != nil {
			responses.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		next(w, r)
	}
}
