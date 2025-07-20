package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func JWTVerifier(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if token == nil {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		t := claims["exp"]
		if exp, ok := t.(time.Time); ok {
			if time.Now().After(exp) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
		} else {
			fmt.Println("Значение не является time.Time")
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Invalid user_id format", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid UUID format", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
