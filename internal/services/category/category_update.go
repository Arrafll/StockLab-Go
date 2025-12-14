package services

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

type CategoryUpdate struct {
	Name string `json:"name" example:"Vitamin"`
}

type CategoryUpdateData struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Vitamin"`
}

type CategoryUpdateSuccessResp struct {
	Status  string             `json:"status" example:"success"`
	Message string             `json:"message" example:"User updated successfully"`
	Data    CategoryUpdateData `json:"data"`
}

type CategoryUpdateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid parameter"`
}

// UpdateCategory godoc
// @Summary Update category product
// @Description Update an existing category product data
// @Tags category
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Category ID"
// @Param name formData string false "Name"
// @Success 200 {object} services.CategoryUpdateSuccessResp
// @Failure 400 {object} services.CategoryUpdateFailResp
// @Failure 404 {object} services.CategoryUpdateFailResp
// @Failure 500 {object} services.CategoryUpdateFailResp
// @Router /stocklab-api/v1/category/update/{id} [put]
// @Security BearerAuth
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "Categgory ID is required")
		return
	}
	catId, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Categgory ID must be a number")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))

	// Cek apakah user ada
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)", catId).Scan(&exists)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if !exists {
		utils.RespondError(w, http.StatusNotFound, "Category not found")
		return
	}

	// Cek name unik (hanya jika name diupdate)
	if name != "" {
		var catExists bool
		err = db.DB.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM categories WHERE LOWER(TRIM(name)) = LOWER($1) AND id <> $2)",
			name, catId,
		).Scan(&catExists)

		if catExists {
			utils.RespondError(w, http.StatusNotFound, "Category with this name exist")
			return
		}
	}

	// Build query dinamis
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if name != "" {
		setParts = append(setParts, "name=$"+strconv.Itoa(argID))
		args = append(args, name)
		argID++
	}

	if len(setParts) == 0 {
		utils.RespondError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	query := "UPDATE categories SET " + strings.Join(setParts, ", ") + " WHERE id=$" + strconv.Itoa(argID) + " RETURNING id, name"
	args = append(args, catId)

	var updatedCategory CategoryUpdateData
	err = db.DB.QueryRow(query, args...).Scan(&updatedCategory.ID, &updatedCategory.Name)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update category: "+err.Error())
		return
	}

	utils.RespondSuccess(w, updatedCategory, "Category updated successfully")
}
