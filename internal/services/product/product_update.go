package services

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
	"github.com/go-chi/chi/v5"
)

// Product blueprint
type ProductUpdateData struct {
	ID         int64  `json:"id" example:"1"`
	Name       string `json:"name" example:"Mie Sedap Goreng"`
	CategoryId *int   `json:"category_id,omitempty" example:"1"`
	SKU        string `json:"sku" example:"SKU-000001"`
	Brand      string `json:"brand" example:"Mie Sedap"`
	Price      string `json:"price" example:"10000"`
	Image      string `json:"image,omitempty" example:"base64imagestring"`
}

type ProductUpdateSuccessResp struct {
	Status  string            `json:"status" example:"success"`
	Message string            `json:"message" example:"Product updated successfully"`
	Data    ProductUpdateData `json:"data"`
}

type ProductUpdateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to update product"`
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product (PATCH semantics)
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param name formData string false "name"
// @Param category_id formData int false "category_id"
// @Param brand formData string false "brand"
// @Param price formData string false "price"
// @Param image formData file false "Product image"
// @Success 200 {object} services.ProductUpdateSuccessResp
// @Failure 400 {object} services.ProductUpdateFailResp
// @Failure 500 {object} services.ProductUpdateFailResp
// @Router /stocklab-api//v1/products/update/{id} [patch]
// @Security BearerAuth
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	name := r.FormValue("name")
	catVal := r.FormValue("category_id")
	brand := r.FormValue("brand")
	price := r.FormValue("price")

	// category_id OPTIONAL
	var categoryID *int
	if catVal != "" {
		val, err := strconv.Atoi(catVal)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "category_id must be number")
			return
		}
		categoryID = &val
	}

	// image OPTIONAL
	var imageBytes []byte
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageBytes, err = io.ReadAll(file)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if err != http.ErrMissingFile {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Build dynamic SET clause
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if name != "" {
		setParts = append(setParts, "name=$"+strconv.Itoa(argID))
		args = append(args, name)
		argID++
	}

	if categoryID != nil {
		setParts = append(setParts, "category_id=$"+strconv.Itoa(argID))
		args = append(args, *categoryID)
		argID++
	}

	if brand != "" {
		setParts = append(setParts, "brand=$"+strconv.Itoa(argID))
		args = append(args, brand)
		argID++
	}

	if price != "" {
		setParts = append(setParts, "price=$"+strconv.Itoa(argID))
		args = append(args, price)
		argID++
	}

	if imageBytes != nil {
		setParts = append(setParts, "image=$"+strconv.Itoa(argID))
		args = append(args, imageBytes)
		argID++
	}

	if len(setParts) == 0 {
		utils.RespondError(w, http.StatusBadRequest, "no fields to update")
		return
	}

	query := `
		UPDATE products
		SET ` + strings.Join(setParts, ", ") + `
		WHERE id=$` + strconv.Itoa(argID) + `
		RETURNING id, sku, name, category_id, brand, price, image
	`

	args = append(args, productID)

	var resp ProductUpdateData
	var imageDB []byte

	err = db.DB.QueryRow(query, args...).Scan(
		&resp.ID,
		&resp.SKU,
		&resp.Name,
		&resp.CategoryId,
		&resp.Brand,
		&resp.Price,
		&imageDB,
	)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if imageDB != nil {
		resp.Image = base64.StdEncoding.EncodeToString(imageDB)
	}

	utils.RespondSuccess(w, resp, "Product updated successfully")
}
