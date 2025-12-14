package services

import (
	"net/http"
	"strconv"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProductDeleteSuccessResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Product deleted successfully"`
}

type ProductDeleteFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to delete product"`
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} services.ProductDeleteSuccessResp
// @Failure 400 {object} services.ProductDeleteFailResp
// @Failure 404 {object} services.ProductDeleteFailResp
// @Failure 500 {object} services.ProductDeleteFailResp
// @Router /stocklab-api/v1/products/delete/{id} [delete]
// @Security BearerAuth
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	query := `
		DELETE FROM products
		WHERE id = $1
		RETURNING id
	`

	var deletedID int64
	err = db.DB.QueryRow(query, productID).Scan(&deletedID)
	if err != nil {
		// Jika ID tidak ditemukan
		if err.Error() == "sql: no rows in result set" {
			utils.RespondError(w, http.StatusNotFound, "product not found")
			return
		}

		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Response sukses
	response := map[string]interface{}{
		"id": productID,
	}

	utils.RespondSuccess(w, response, "Product deleted successfully")
}
