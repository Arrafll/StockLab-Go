package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID   int    `json:"id" example:"1"` // ID kapital
	Name string `json:"name" example:"Andre"`
}

type Resp struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// GetUserList godoc
// @Summary Get list of users
// @Description Get all users in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} services.User
// @Router /api/v1/users/ [get]
// @Security BearerAuth
func GetUserList(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "Andre"},
		{ID: 2, Name: "Rafli"},
		{ID: 3, Name: "Brother"},
	}

	resp := Resp{
		Status: "success",
		Data:   users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DetailUser godoc
// @Summary Get detail of an user
// @Description Get all users in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} services.User
// @Router /api/v1/users/detail/{id} [get]
func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "Andre"},
		{ID: 2, Name: "Rafli"},
		{ID: 3, Name: "Brother"},
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	var userDetail *User
	for i := range users {
		if users[i].ID == idInt {
			userDetail = &users[i]
			break
		}
	}

	if userDetail == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Println("User Detail:", userDetail)

	resp := Resp{
		Status: "success",
		Data:   userDetail,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
