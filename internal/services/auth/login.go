package services

import (
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" example:"andrerafli83@gmail.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJI..."`
}

type LoginFailResponse struct {
	Message string `json:"message" example:"Invalid credentials"`
}

// Login godoc
// @Summary User login
// @Description Login to the system
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} services.LoginResponse
// @Failure 400 {object} map[string]string "Request invalid"
// @Param body body services.LoginRequest true "User login credentials"
// @Router /api/v1/login/ [post]
func Login(w http.ResponseWriter, r *http.Request) {
	// Implementation for user login

}
