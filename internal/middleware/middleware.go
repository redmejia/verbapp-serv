package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redmejia/internal/handlers"
	"github.com/redmejia/internal/models"
	"github.com/redmejia/internal/security"
)

func IsAuthorized(app *handlers.App, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		autorization := r.Header.Get("Authorization")

		if len(autorization) > 0 {
			token := strings.Split(autorization, " ")
			isValid, _, err := security.VerifyToken(token[1], app.JwtKey)

			if err != nil {

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				if errors.Is(err, jwt.ErrTokenExpired) {
					json.NewEncoder(w).Encode(models.ErrorResponse{
						Code:    http.StatusUnauthorized,
						Message: "Token expired",
					})
				} else if errors.Is(err, jwt.ErrTokenMalformed) {
					json.NewEncoder(w).Encode(models.ErrorResponse{
						Code:    http.StatusUnauthorized,
						Message: "Token malformed",
					})
				} else {
					json.NewEncoder(w).Encode(models.ErrorResponse{
						Code:    http.StatusUnauthorized,
						Message: "Unauthorized",
					})
				}
				return
			}
			if isValid {
				next.ServeHTTP(w, r)
			}
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
		}
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.Proto)
		next.ServeHTTP(w, r)
		log.Println("Request took: ", time.Since(start))
	})
}
