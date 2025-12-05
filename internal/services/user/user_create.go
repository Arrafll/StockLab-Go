package services

import "net/http"

type UserCreate struct {
	Name     string `json:"name" example:"Andre"`
	Email    string `json:"email" example:andrerafli83@gmail.com"`
	Password string `json:"password" example:"password123"`
}

// CreateUser godoc
// @Summary Get detail of an user
// @Description Get all users in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} services.User
// @Router /api/v1/users/create/ [post]
// @Param user body UserCreate true "Data pengguna"
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a user
}
