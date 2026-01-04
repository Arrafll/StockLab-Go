package services

import (
	"net/http"
	"strconv"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

// CategoryDetailSuccessResp untuk response detail user
type CategoryDeleteSuccessResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Category deleted successfully"`
}

// CategoryDeleteFailResp untuk error
type CategoryDeleteFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to delete category"`
}

// DeleteCategory godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {object} services.CategoryDeleteSuccessResp
// @Failure 400 {object} services.CategoryDeleteFailResp
// @Failure 500 {object} services.CategoryDeleteFailResp
// @Param id path int true "Category Id"
// @Router /stocklab-api/v1/categories/delete/{id} [delete]
// @Security BearerAuth
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "Category Id is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Category Id must be a number")
		return
	}

	// Cek apakah user ada
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if !exists {
		utils.RespondError(w, http.StatusNotFound, "Category not found")
		return
	}

	// Cek apakah ada product dengan category ini
	var prodExist bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE category_id=$1)", id).Scan(&prodExist)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if prodExist {
		utils.RespondError(w, http.StatusInternalServerError, "Cannot delete category, products with this category exist")
		return
	}

	// Delete user
	_, err = db.DB.Exec("DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete category: "+err.Error())
		return
	}

	// Response sukses
	response := map[string]interface{}{
		"id": id,
	}
	utils.RespondSuccess(w, response, "Category deleted successfully")
}
