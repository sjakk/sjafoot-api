package main


import (
	"context"
	"errors"
	"net/http"
	"github.com/sjakk/sjafoot/internal/data"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const userContextKey = contextKey("user")

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.errorResponse(w, r, http.StatusUnauthorized, "Authorization header is missing")
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.errorResponse(w, r, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := headerParts[1]
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(app.config.jwt.secret), nil
		})

		if err != nil || !token.Valid {
			app.errorResponse(w, r, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		userIDFloat, ok := (*claims)["sub"].(float64)
		if !ok {
			app.serverErrorResponse(w, r, errors.New("invalid user ID in token"))
			return
		}
		userID := int64(userIDFloat)

		user, err := app.models.Users.GetByID(userID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Add the user to the request context.
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAdminUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(*data.User)
		if !ok {
			app.serverErrorResponse(w, r, errors.New("user not found in context"))
			return
		}

		if user.Role != "admin" {
			app.errorResponse(w, r, http.StatusForbidden, "You do not have the necessary permissions to access this resource")
			return
		}

		next.ServeHTTP(w, r)
	})
}


