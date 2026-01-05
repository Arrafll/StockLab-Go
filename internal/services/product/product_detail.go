package services

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

// Product Detail blueprint
type ProductDetail struct {
	ID       int64  `json:"id" example:"1"`
	Name     string `json:"name" example:"Mie Sedap Goreng"`
	Category string `json:"category" example:"Mie"`
	SKU      string `json:"sku" example:"SKU-20251214201530-042"`
	Brand    string `json:"brand" example:"Mie Sedap"`
	Price    string `json:"price" example:"10000"`
	Quantity int32  `json:"quantity" example:"150"`
	Image    string `json:"image" example:"base64imagestring"`
}

type ProductDetailSuccessResp struct {
	Status  string        `json:"status" example:"success"`
	Message string        `json:"message" example:"Product fetched successfully"`
	Data    ProductDetail `json:"data"`
}

type ProductDetailFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch product"`
}

// GetProductDetail godoc
// @Summary Product detail
// @Description Menampilkan detail product berdasarkan ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} services.ProductDetailSuccessResp
// @Failure 400 {object} services.ProductDetailFailResp
// @Failure 404 {object} services.ProductDetailFailResp
// @Failure 500 {object} services.ProductDetailFailResp
// @Router /stocklab-api/v1/products/detail/{id} [get]
// @Security BearerAuth
func GetProductDetail(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	idParam := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	query := `
		SELECT 
			p.id,
			p.name,
			p.sku,
			p.brand,
			COALESCE(p.price, '0') as price,
			COALESCE(c.name, 'N/A') AS category,
			COALESCE(s.quantity, 0) as quantity,
			p.image
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		LEFT JOIN stocks s ON s.product_id = p.id
		WHERE p.id = $1
	`

	var (
		imageBytes []byte
		product    ProductDetail
	)

	err = db.DB.QueryRow(query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Brand,
		&product.Price,
		&product.Category,
		&product.Quantity,
		&imageBytes,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.RespondError(w, http.StatusNotFound, "product not found")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert image ke base64 jika ada
	if imageBytes != nil {
		product.Image = base64.StdEncoding.EncodeToString(imageBytes)
	} else {
		product.Image = ""
	}

	utils.RespondSuccess(w, product, "Product fetched successfully")
}
