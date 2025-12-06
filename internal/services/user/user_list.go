package services

import (
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// User model
type User struct {
	ID    int    `json:"id" example:"1"`
	Email string `json:"email" example:"andrerafli83@gmail.com"`
	Name  string `json:"name" example:"Andre"`
	Phone string `json:"phone" example:"081234567890"`
	Role  string `json:"role" example:"staff"`
}

// GetUserList godoc
// @Summary Get list of users
// @Description Get all users in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} Resp
// @Router /stocklab-api/v1/users/list [get]
// @Security BearerAuth
func GetUserList(w http.ResponseWriter, r *http.Request) {
	// Query semua user
	rows, err := db.DB.Query("SELECT id, email, name, phone, role FROM users ORDER BY id ASC")
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.Role); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to scan user: "+err.Error())
			return
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
