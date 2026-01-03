package services

import (
	"encoding/base64"
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// Product blueprint
type Product struct {
	ID       int    `json:"id" example:"1"`
	Name     string `json:"name" example:"Mie Sedap Goreng"`
	Category string `json:"category" example:"Mie"`
	SKU      string `json:"sku" example:"SKU-20251214201530-042"`
	Brand    string `json:"brand" example:"Mie Sedap"`
	Image    string `json:"image" form:"image" example:"base64imagestring"`
}

type ProductSuccessResp struct {
	Status  string    `json:"status" example:"success"`
	Message string    `json:"message" example:"Product fetched successfully"`
	Data    []Product `json:"data"`
}

type ProductFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch product"`
}

// Product godoc
// @Summary Product list
// @Description Menampilkan list product
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {object} services.ProductSuccessResp
// @Failure 500 {object} services.ProductFailResp
// @Security BearerAuth
// @Router /stocklab-api//v1/products/ [get]
func GetProductList(w http.ResponseWriter, r *http.Request) {
	// Query semua category
	rows, err := db.DB.Query("SELECT p.id, p.name, p.sku, p.brand, COALESCE( c.name, 'N/A') as category, p.image FROM products p LEFT JOIN categories c ON c.id = p.category_id ORDER BY id DESC")
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch products: "+err.Error())
		return
	}
	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var image []byte
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.SKU, &p.Brand, &image); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to scan products: "+err.Error())
			return
		}

		//Condition null image
		if image != nil {
			p.Image = base64.StdEncoding.EncodeToString(image)
		} else {
			p.Image = ""
		}

		products = append(products, p)
	}

	// Cek apakah ada error saat iterasi rows
	if err = rows.Err(); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error reading products: "+err.Error())
		return
	}

	utils.RespondSuccess(w, products, "Products fetched successfully")
}
