package services

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

// User model
type UserDetail struct {
	ID     int       `json:"id" example:"1"`
	Email  string    `json:"email" example:"andrerafli83@gmail.com"`
	Name   string    `json:"name" example:"Andre"`
	Phone  string    `json:"phone" example:"081234567890"`
	Role   string    `json:"role" example:"staff"`
	Joined time.Time `json:"joined" example:"2025-12-14"`
	Avatar string    `json:"avatar"example:"base64imagestring"`
}

// UserDetailSuccessResp untuk response detail user
type UserDetailSuccessResp struct {
	Status  string     `json:"status" example:"success"`
	Message string     `json:"message" example:"User fetched successfully"`
	Data    UserDetail `json:"data"`
}

// UserDetailFailResp untuk error
type UserDetailFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch user"`
}

// GetUserDetail godoc
// @Summary Get detail of a user
// @Description Get a single user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} services.UserDetailSuccessResp
// @Failure 404 {object} services.UserDetailFailResp
// @Failure 500 {object} services.UserDetailFailResp
// @Param id path int true "User ID"
// @Router /stocklab-api/v1/users/detail/{id} [get]
// @Security BearerAuth
func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "User ID must be a number")
		return
	}

	// Query user by ID
	var u UserDetail
	var avatar []byte
	var joined time.Time
	query := `SELECT id, email, name, phone, role, avatar, created_at FROM users WHERE id = $1`
	err = db.DB.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.Role, &avatar, &joined)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		return
	}

	if avatar != nil {
		u.Avatar = base64.StdEncoding.EncodeToString(avatar)
	} else {
		u.Avatar = ""
	}

	u.Joined = joined

	// Response sukses
	resp := UserDetailSuccessResp{
		Status:  "success",
		Message: "User fetched successfully",
		Data:    u,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
