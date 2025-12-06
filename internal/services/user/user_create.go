package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserCreate struct {
	Email    string `json:"email" example:andrerafli83@gmail.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"Andre"`
	Phone    string `json:"phone" example:"09999999999"`
}

type UserCreateData struct {
	ID    int64  `json:"id" example:"1"`
	Email string `json:"email" example:"andrerafli83@gmail.com"`
	Name  string `json:"name" example:"Andre"`
	Phone string `json:"phone" example:"09999999999"`
	Role  string `json:"role" example:"staff"`
}

type UserCreateSuccessResp struct {
	Status  string         `json:"status" example:"success"`
	Message string         `json:"message" example:"Login successful"`
	Data    UserCreateData `json:"data"`
}

type UserCreateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid credentials"`
}

// CreateUser godoc
// @Summary Get detail of an user
// @Description Get all users in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} services.UserCreateSuccessResp
// @Failure 400 {object} services.UserCreateFailResp
// @Router /stocklab-api/v1/users/create [post]
// @Param user body UserCreate true "Data pengguna"
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser UserCreate
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if newUser.Email == "" || newUser.Password == "" || newUser.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email, password, and name are required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	// Cek apakah email sudah ada
	newUser.Email = strings.TrimSpace(newUser.Email)

	newUser.Email = strings.TrimSpace(newUser.Email)

	var exists bool
	err = db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(TRIM(email)) = LOWER($1))",
		newUser.Email,
	).Scan(&exists)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	if exists {
		utils.RespondError(w, http.StatusConflict, newUser.Email+" is already registered")
		return
	}

	// Insert user ke database
	query := `INSERT INTO users (email, password, name, phone, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var userID int64
	err = db.DB.QueryRow(query, newUser.Email, string(hashedPassword), newUser.Name, newUser.Phone, "staff").Scan(&userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	// Response sukses
	response := map[string]interface{}{
		"id":    userID,
		"email": newUser.Email,
		"name":  newUser.Name,
		"phone": newUser.Phone,
	}

	utils.RespondSuccess(w, response, "User created successfully")
}
