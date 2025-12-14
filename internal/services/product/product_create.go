package services

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// Product blueprint
type ProductCreateData struct {
	ID         int64  `json:"id" example:"1"`
	Name       string `json:"name" example:"Mie Sedap Goreng"`
	CategoryId int    `json:"category_id" example:"1"`
	SKU        string `json:"sku" example:"SKU-20251214201530-042"`
	Brand      string `json:"brand" example:"Mie Sedap"`
	Price      string `json:"price" example:"10000"`
	Image      string `json:"image" form:"image" example:"base64imagestring"`
}

type ProductCreateSuccessResp struct {
	Status  string              `json:"status" example:"success"`
	Message string              `json:"message" example:"Product created successfully"`
	Data    []ProductCreateData `json:"data"`
}

type ProductCreateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to create product"`
}

// CreateCategory godoc
// @Summary Create product for products
// @Description Create a product
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "name"
// @Param category_id formData int true "category_id"
// @Param brand formData string true "brand"
// @Param image formData file true "Product image"
// @Success 200 {object} services.ProductCreateData
// @Failure 400 {object} services.ProductCreateFailResp
// @Failure 500 {object} services.ProductCreateFailResp
// @Router /stocklab-api/v1/products/create [post]
// @Security BearerAuth
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
		return
	}

	name := r.FormValue("name")
	catId := r.FormValue("category_id")
	brand := r.FormValue("brand")
	price := r.FormValue("price")

	categoryId, err := strconv.Atoi(catId)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Category Id must be a number")
		return
	}

	// Ambil file avatar
	file, _, err := r.FormFile("image")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Failed to read image: "+err.Error())
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to read image bytes: "+err.Error())
		return
	}

	sku := GenerateSKU()

	// Insert product ke database
	var productId int64
	query := `INSERT INTO products (name, category_id, sku, brand, price, image) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = db.DB.QueryRow(query, name, categoryId, sku, brand, price, imageBytes).Scan(&productId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create product: "+err.Error())
		return
	}

	// Response sukses
	response := ProductCreateData{
		ID:         productId,
		Name:       name,
		CategoryId: categoryId,
		SKU:        sku,
		Brand:      brand,
		Price:      price,
		Image:      base64.StdEncoding.EncodeToString(imageBytes),
	}

	utils.RespondSuccess(w, response, "Product created successfully")
}

// Simple generate SKU
func GenerateSKU() string {
	ts := time.Now().UTC().Format("20060102150405") // yyyyMMddHHmmss
	randPart := rand.Intn(1000)                     // 000–999
	return fmt.Sprintf("SKU-%s-%03d", ts, randPart)
}
