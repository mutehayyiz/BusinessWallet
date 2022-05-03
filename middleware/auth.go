package middleware

import (
	"BusinessWallet/auth"
	"BusinessWallet/controller"
	"context"
	"net/http"
	"strings"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			controller.ReturnError(w, http.StatusUnauthorized, "Malformed token")
			return
		}
		jwtToken := authHeader[1]
		token, _ := auth.VerifyToken(jwtToken)
		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next(w, r.WithContext(ctx))
		} else {
			controller.ReturnError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
	}
}
