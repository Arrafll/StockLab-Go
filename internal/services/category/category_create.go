package services

import (
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

type CategoryCreate struct {
	Name string `json:"name" example:"Vitamin"`
}

type CategoryCreateData struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Vitamin"`
}

type CategoryCreateSuccessResp struct {
	Status  string             `json:"status" example:"success"`
	Message string             `json:"message" example:"Category created successfully"`
	Data    CategoryCreateData `json:"data"`
}

type CategoryCreateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid parameter"`
}

// CreateCategory godoc
// @Summary Create category for products
// @Description Create a category
// @Tags category
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Success 200 {object} services.CategoryCreateSuccessResp
// @Failure 400 {object} services.CategoryCreateFailResp
// @Failure 500 {object} services.CategoryCreateFailResp
// @Router /stocklab-api/v1/category/create [post]
// @Security BearerAuth
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
		return
	}

	name := r.FormValue("name")

	if name == "" {
		utils.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}

	// Cek apakah name sudah ada
	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM categories WHERE LOWER(TRIM(name)) = LOWER($1))",
		name,
	).Scan(&exists)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, name+" is already registered")
		return
	}
	// Insert category ke database
	var catId int64
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	err = db.DB.QueryRow(query, name).Scan(&catId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create category: "+err.Error())
		return
	}

	// Response sukses
	response := CategoryCreateData{
		ID:   catId,
		Name: name,
	}

	utils.RespondSuccess(w, response, "Category created successfully")
}
