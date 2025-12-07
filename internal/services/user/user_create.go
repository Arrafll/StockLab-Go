package services

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserCreate struct {
	Email    string `json:"email" example:"andrerafli83@gmail.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"Andre"`
	Phone    string `json:"phone" example:"09999999999"`
	Avatar   string `json:"avatar" form:"avatar" example:"file"`
}

type UserCreateData struct {
	ID     int64  `json:"id" example:"1"`
	Email  string `json:"email" example:"andrerafli83@gmail.com"`
	Name   string `json:"name" example:"Andre"`
	Phone  string `json:"phone" example:"09999999999"`
	Role   string `json:"role" example:"staff"`
	Avatar string `json:"avatar" form:"avatar" example:"base64imagestring"`
}

type UserCreateSuccessResp struct {
	Status  string         `json:"status" example:"success"`
	Message string         `json:"message" example:"User created successfully"`
	Data    UserCreateData `json:"data"`
}

type UserCreateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid credentials"`
}

// CreateUser godoc
// @Summary Create user with avatar
// @Description Create a new user and upload avatar
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string true "Name"
// @Param phone formData string false "Phone"
// @Param avatar formData file true "Avatar image"
// @Success 200 {object} services.UserCreateSuccessResp
// @Failure 400 {object} services.UserCreateFailResp
// @Failure 500 {object} services.UserCreateFailResp
// @Router /stocklab-api/v1/users/create [post]
// @Security BearerAuth
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	if email == "" || password == "" || name == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email, password, and name are required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Cek apakah email sudah ada
	var exists bool
	err = db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(TRIM(email)) = LOWER($1))",
		email,
	).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, email+" is already registered")
		return
	}

	// Ambil file avatar
	file, _, err := r.FormFile("avatar")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Failed to read avatar: "+err.Error())
		return
	}
	defer file.Close()

	avatarBytes, err := io.ReadAll(file)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to read avatar bytes: "+err.Error())
		return
	}

	// Insert user ke database
	var userID int64
	query := `INSERT INTO users (email, password, name, phone, role, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = db.DB.QueryRow(query, email, string(hashedPassword), name, phone, "staff", avatarBytes).Scan(&userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	// Response sukses
	response := UserCreateData{
		ID:     userID,
		Email:  email,
		Name:   name,
		Phone:  phone,
		Role:   "staff",
		Avatar: base64.StdEncoding.EncodeToString(avatarBytes),
	}

	utils.RespondSuccess(w, response, "User created successfully")
}
