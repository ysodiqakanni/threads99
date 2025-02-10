package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ysodiqakanni/threads99/internal/helper"
	"net/http"
)

func AuthenticateMiddleware(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your authentication logic here
		// Check if the JWT token is valid and extract user information
		// For example, you can check the "Authorization" header for the token
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			//http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			helper.EncodeErrorResponse(w, errors.New("Authorization token missing"),
				"Authorization token missing", "401")
			return
		}

		// test token validation
		var tokenSecret = []byte(jwtSecret)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method and return the secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid token signing method")
			}
			return tokenSecret, nil
		})

		if err != nil || !token.Valid {
			fmt.Println(err.Error())
			helper.EncodeErrorResponse(w, errors.New("Login expired"),
				"Your session has expired. Please login again.", "401")
			return
		}
		// end test token validation. Let's inspect the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helper.EncodeErrorResponse(w, errors.New("Invalid token"),
				"Your session has expired. Please login again.", "401")
			return
		}
		// the jwt library internally converts my jwt []string to []interface.
		// So Ima convert to []string for easy manipulation

		/*
			roles, ok := claims["role"].([]interface{})
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Convert []interface{} to []string
			var rolesSlice []string
			for _, role := range roles {
				rolesSlice = append(rolesSlice, fmt.Sprintf("%v", role))
			}
		*/

		fmt.Println("User claims are: ", claims)
		// Add user information to the request context
		ctx := context.WithValue(r.Context(), "username", claims["username"].(string))
		ctx = context.WithValue(ctx, "userId", claims["id"].(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
