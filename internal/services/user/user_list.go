package services

import (
	"encoding/base64"
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// User model
type User struct {
	ID     int    `json:"id" example:"1"`
	Email  string `json:"email" example:"andrerafli83@gmail.com"`
	Name   string `json:"name" example:"Andre"`
	Phone  string `json:"phone" example:"081234567890"`
	Role   string `json:"role" example:"staff"`
	Avatar string `json:"avatar" form:"avatar" example:"base64imagestring"`
}

type UserListSuccessResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Login successful"`
	Data    []User `json:"data"`
}

type UserListFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch users"`
}

// GetUserList godoc
// @Summary Get list of users
// @Description Get all users in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} services.UserListSuccessResp
// @Failure 500 {object} services.UserListFailResp
// @Router /stocklab-api/v1/users [get]
// @Security BearerAuth
func GetUserList(w http.ResponseWriter, r *http.Request) {
	// Query semua user
	rows, err := db.DB.Query("SELECT id, email, name, phone, role, avatar FROM users ORDER BY id ASC")
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		var avatar []byte
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.Role, &avatar); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to scan user: "+err.Error())
			return
		}

		//Condition null image
		if avatar != nil {
			u.Avatar = base64.StdEncoding.EncodeToString(avatar)
		} else {
			u.Avatar = ""
		}
		users = append(users, u)
	}

	// Cek apakah ada error saat iterasi rows
	if err = rows.Err(); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error reading users: "+err.Error())
		return
	}

	utils.RespondSuccess(w, users, "Users fetched successfully")
}
