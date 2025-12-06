package services

import (
	"context"
	"net/http"
	"strings"

	"github.com/Arrafll/StockLab-Go/internal/config"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// JWTMiddleware membuat middleware untuk memvalidasi token
func JWTMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil token dari header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}

			tokenStr := parts[1]

			// Validasi token
			claims, err := utils.ValidateJWT(tokenStr, cfg.JWTSecret)
			if err != nil {
				utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}

			// Simpan user_id di context agar bisa dipakai handler
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims["user_id"])
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
