package services

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdate struct {
	Email    string `json:"email" example:"andrerafli83@gmail.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"Andre"`
	Phone    string `json:"phone" example:"09999999999"`
	Avatar   string `json:"avatar" form:"avatar" example:"file"`
}

type UserUpdateData struct {
	ID     int64  `json:"id" example:"1"`
	Email  string `json:"email" example:"andrerafli83@gmail.com"`
	Name   string `json:"name" example:"Andre"`
	Phone  string `json:"phone" example:"09999999999"`
	Role   string `json:"role" example:"staff"`
	Avatar string `json:"avatar" form:"avatar" example:"base64imagestring"`
}

type UserUpdateSuccessResp struct {
	Status  string         `json:"status" example:"success"`
	Message string         `json:"message" example:"User updated successfully"`
	Data    UserUpdateData `json:"data"`
}

type UserUpdateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid parameter"`
}

// UpdateUser godoc
// @Summary Update user with avatar
// @Description Update an existing user and optionally upload a new avatar
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "User ID"
// @Param email formData string false "Email"
// @Param password formData string false "Password"
// @Param name formData string false "Name"
// @Param phone formData string false "Phone"
// @Param avatar formData file false "Avatar image"
// @Success 200 {object} services.UserUpdateSuccessResp
// @Failure 400 {object} services.UserUpdateFailResp
// @Failure 404 {object} services.UserUpdateFailResp
// @Failure 500 {object} services.UserUpdateFailResp
// @Router /stocklab-api/v1/users/update/{id} [put]
// @Security BearerAuth
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "User ID is required")
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "User ID must be a number")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	// Ambil file avatar jika ada
	var avatarBytes []byte
	file, _, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarBytes, err = io.ReadAll(file)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to read avatar bytes: "+err.Error())
			return
		}
	}

	// Cek apakah user ada
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", userID).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if !exists {
		utils.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Cek email unik (hanya jika email diupdate)
	if email != "" {
		var emailExists bool
		err = db.DB.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(TRIM(email)) = LOWER($1) AND id <> $2)",
			email, userID,
		).Scan(&emailExists)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}
		if emailExists {
			utils.RespondError(w, http.StatusConflict, email+" is already registered by another user")
			return
		}
	}

	// Build query dinamis
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if email != "" {
		setParts = append(setParts, "email=$"+strconv.Itoa(argID))
		args = append(args, email)
		argID++
	}
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		setParts = append(setParts, "password=$"+strconv.Itoa(argID))
		args = append(args, string(hashedPassword))
		argID++
	}
	if name != "" {
		setParts = append(setParts, "name=$"+strconv.Itoa(argID))
		args = append(args, name)
		argID++
	}
	if phone != "" {
		setParts = append(setParts, "phone=$"+strconv.Itoa(argID))
		args = append(args, phone)
		argID++
	}
	if avatarBytes != nil {
		setParts = append(setParts, "avatar=$"+strconv.Itoa(argID))
		args = append(args, avatarBytes)
		argID++
	}

	if len(setParts) == 0 {
		utils.RespondError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE id=$" + strconv.Itoa(argID) + " RETURNING id, email, name, phone, role, avatar"
	args = append(args, userID)

	var updatedUser UserCreateData
	var avatarDB []byte
	err = db.DB.QueryRow(query, args...).Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.Name, &updatedUser.Phone, &updatedUser.Role, &avatarDB)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	// Convert avatar to base64 string
	if avatarDB != nil {
		updatedUser.Avatar = base64.StdEncoding.EncodeToString(avatarDB)
	}

	utils.RespondSuccess(w, updatedUser, "User updated successfully")
}
