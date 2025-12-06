package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Arrafll/StockLab-Go/internal/config"
	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginParamRequest struct {
	Email    string `json:"email" example:"andrerafli83@gmail.com"`
	Password string `json:"password" example:"password123"`
}

type AuthLoginData struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AuthLoginSuccessResponse struct {
	Status  string        `json:"status" example:"success"`
	Message string        `json:"message" example:"Login successful"`
	Data    AuthLoginData `json:"data"`
}

type AuthLoginFailResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid credentials"`
}

// Login godoc
// @Summary User login
// @Description Login to the system
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} services.AuthLoginSuccessResponse
// @Failure 401 {object} services.AuthLoginFailResponse
// @Param body body services.AuthLoginParamRequest true "User login credentials"
// @Router /stocklab-api/v1/login [post]
func Login(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	var req AuthLoginParamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var storedHash string
	var userID int64
	err := db.DB.QueryRowContext(ctx, `SELECT id,password FROM users WHERE email=$1 LIMIT 1`, req.Email).
		Scan(&userID, &storedHash)

	if err == sql.ErrNoRows {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenStr, err := utils.GenerateJWT(userID, cfg.JWTSecret)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	utils.RespondSuccess(w, map[string]string{"token": tokenStr}, "Login successful")
}
