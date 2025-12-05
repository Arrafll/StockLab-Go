package services

import "net/http"

type ProductResponse struct {
	ProductID   int     `json:"product_id" example:"101"`
	ProductName string  `json:"product_name" example:"Sample Product"`
	Price       float64 `json:"price" example:"29.99"`
}

// Product godoc
// @Summary Product list
// @Description Menampilkan list product
// @Tags product
// @Accept  json
// @Produce  json
// @Success 200 {object} services.ProductResponse
// @Router /api/v1/product/ [post]
func GetProducts(w http.ResponseWriter, r *http.Request) {

}
