package services

import (
	"net/http"
	"time"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

type TransactionListData struct {
	ID           int64     `json:"id" example:"1"`
	ProductName  string    `json:"product_name" example:"Mie Sedap Goreng"`
	ProductSKU   string    `json:"product_sku" example:"MIE001"`
	ProductBrand string    `json:"product_brand" example:"Sedap"`
	ProductPrice int64     `json:"product_price" example:"5000"`
	PICName      string    `json:"pic_name" example:"John Doe"`
	Quantity     int64     `json:"quantity" example:"10"`
	MoveType     string    `json:"move_type" example:"in"`
	CreatedAt    time.Time `json:"created_at" example:"2024-12-14T20:15:30Z"` // ISO 8601 format
}
type TransactionListSuccessResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Transaction fetched successfully"`
}

type TransactionListFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch transaction"`
}

// ListTransaction godoc
// @Summary List transaction stocks
// @Description List a transaction for stock movements
// @Tags transactions
// @Accept multipart/form-data
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} services.TransactionListData
// @Failure 400 {object} services.TransactionListFailResp
// @Failure 500 {object} services.TransactionListFailResp
// @Router /stocklab-api/v1/transactions/list [post]
// @Security BearerAuth
func GetTransactionList(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	query := `
		SELECT tr.id, p.name as product_name, p.sku, p.brand, COALESCE(p.price, 0) as price, u.name as pic_name, tr.quantity, tr.move_type, tr.created_at
		FROM transactions tr
		LEFT JOIN products p ON tr.product_id = p.id
		LEFT JOIN users u ON tr.user_id = u.id
	`

	var args []interface{}
	where := ""

	// Build WHERE condition
	if startDate != "" && endDate != "" {
		where = "WHERE tr.created_at::date BETWEEN $1 AND $2"
		args = append(args, startDate, endDate)
	} else if startDate != "" {
		where = "WHERE tr.created_at::date >= $1"
		args = append(args, startDate)
	} else if endDate != "" {
		where = "WHERE tr.created_at::date <= $1"
		args = append(args, endDate)
	}

	order := " ORDER BY tr.created_at DESC"

	rows, err := db.DB.Query(query+where+order, args...)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed fetching transactions: "+err.Error())
		return
	}
	defer rows.Close()

	var data []TransactionListData

	for rows.Next() {
		var t TransactionListData
		if err := rows.Scan(
			&t.ID,
			&t.ProductName,
			&t.ProductSKU,
			&t.ProductBrand,
			&t.ProductPrice,
			&t.PICName,
			&t.Quantity,
			&t.MoveType,
			&t.CreatedAt,
		); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed parsing transactions: "+err.Error())
			return
		}

		data = append(data, t)
	}

	utils.RespondSuccess(w, data, "Transactions fetched successfully")
}
