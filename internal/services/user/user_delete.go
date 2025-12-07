package services

import (
	"net/http"
	"strconv"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

// UserDetailSuccessResp untuk response detail user
type UserDeleteSuccessResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"User deleted successfully"`
}

// UserDeleteFailResp untuk error
type UserDeleteFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to delete user"`
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} services.UserDeleteSuccessResp
// @Failure 400 {object} services.UserDeleteFailResp
// @Failure 500 {object} services.UserDeleteFailResp
// @Param id path int true "User ID"
// @Router /stocklab-api/v1/users/delete/{id} [delete]
// @Security BearerAuth
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
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

	// Cek apakah user ada
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if !exists {
		utils.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Delete user
	_, err = db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	// Response sukses
	response := map[string]interface{}{
		"id": id,
	}
	utils.RespondSuccess(w, response, "User deleted successfully")
}
