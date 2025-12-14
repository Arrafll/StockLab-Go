package services

import (
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// Category blueprint
type Category struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Vitamin"`
}

type CategorySuccessResp struct {
	Status  string     `json:"status" example:"success"`
	Message string     `json:"message" example:"Categories fetched successfully"`
	Data    []Category `json:"data"`
}

type CategoryFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch categories"`
}

// GetCategoryList godoc
// @Summary Get list of category
// @Description Get all category in the system
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {object} services.CategorySuccessResp
// @Failure 500 {object} services.CategoryFailResp
// @Router /stocklab-api/v1/categories [get]
// @Security BearerAuth
func GetCategoryList(w http.ResponseWriter, r *http.Request) {
	// Query semua category
	rows, err := db.DB.Query("SELECT id,name FROM categories ORDER BY id DESC")
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch categories: "+err.Error())
		return
	}
	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to scan categories: "+err.Error())
			return
		}

		categories = append(categories, c)
	}

	// Cek apakah ada error saat iterasi rows
	if err = rows.Err(); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error reading categories: "+err.Error())
		return
	}

	utils.RespondSuccess(w, categories, "Categories fetched successfully")
}
