package services

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// CreateTransaction godoc
// @Summary Create transaction stocks
// @Description Create a transaction for stock movements
// @Tags transactions
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData int true "product_id"
// @Param user_id formData int true "user_id"
// @Param quantity formData int true "quantity"
// @Param move_type formData string true "move_type"
// @Success 200 {object} services.TransactionCreateData
// @Failure 400 {object} services.TransactionCreateFailResp
// @Failure 500 {object} services.TransactionCreateFailResp
// @Router /stocklab-api/v1/transactions/create [post]
// @Security BearerAuth
type TransactionCreateData struct {
	ID        int64  `json:"id" example:"1"`
	ProductID int64  `json:"product_id" example:"1"`
	UserID    int64  `json:"user_id"  example:"1"`
	Quantity  int64  `json:"quantity" example:"100"`
	MoveType  string `json:"move_type" example:"in"` // IN | OUT
}

type TransactionCreateSuccessResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Transaction created successfully"`
}

type TransactionCreateFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to create transaction"`
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	productID, _ := strconv.ParseInt(r.FormValue("product_id"), 10, 64)
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	qty, _ := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	moveType := strings.ToUpper(r.FormValue("move_type"))

	if productID == 0 || qty <= 0 || (moveType != "IN" && moveType != "OUT") {
		utils.RespondError(w, http.StatusBadRequest, "Invalid transaction payload")
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to start transaction: "+err.Error())
		return
	}
	defer tx.Rollback()

	var currentQty int64

	// Lock stock row to prevent race conditions
	err = tx.QueryRow(`
        SELECT quantity
        FROM stocks
        WHERE product_id = $1
        FOR UPDATE
    `, productID).Scan(&currentQty)

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Stock not found: "+err.Error())
		return
	}

	// Apply movement
	if moveType == "OUT" {
		if currentQty-qty < 0 {
			utils.RespondError(w, http.StatusConflict, "Insufficient stock")
			return
		}
		currentQty -= qty
	} else {
		currentQty += qty
	}

	// Update stock
	_, err = tx.Exec(`
        UPDATE stocks
        SET quantity = $1, updated_at = $2
        WHERE product_id = $3
    `, currentQty, time.Now(), productID)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed updating stock: "+err.Error())
		return
	}

	// Insert transaction history
	var txId int64
	err = tx.QueryRow(`
        INSERT INTO transactions (product_id, user_id, quantity, move_type)
        VALUES ($1, $2, $3, $4)
		RETURNING id
    `, productID, userID, qty, moveType).Scan(&txId)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed creating transaction: "+err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Commit failed: "+err.Error())
		return
	}

	utils.RespondSuccess(
		w,
		TransactionCreateData{
			ID:        txId,
			ProductID: productID,
			Quantity:  qty,
			MoveType:  moveType,
		},
		"Transaction created successfully",
	)
}
