package services

import (
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// Dashboard blueprint
type DashboardData struct {
	ProductTotal      int                      `json:"product_total" example:"100"`
	StockTotal        int                      `json:"stock_total" example:"5000"`
	LowStockTotal     int                      `json:"low_stock" example:"50"`
	NoStockTotal      int                      `json:"no_stock" example:"10"`
	ChartActivityData []map[string]interface{} `json:"chart_activity_data"`
}

type DashboardSuccessResp struct {
	Status  string        `json:"status" example:"success"`
	Message string        `json:"message" example:"Dashboard data fetched successfully"`
	Data    DashboardData `json:"data"`
}

type DashboardFailResp struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed to fetch dashboard data"`
}

// Dashboard godoc
// @Summary Dashboard data
// @Description Menampilkan data dashboard
// @Tags dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object} services.DashboardSuccessResp
// @Failure 500 {object} services.DashboardFailResp
// @Security BearerAuth
// @Router /stocklab-api/v1/dashboard/ [get]
func DashboardMain(w http.ResponseWriter, r *http.Request) {
	// Query semua category
	var DashboardData DashboardData
	widgetQuery := `SELECT 
		(SELECT COUNT(*) FROM products) as product_total,
		(SELECT SUM(quantity) FROM stocks) as stock_total,
		(SELECT COUNT(*) FROM stocks WHERE quantity < 10 AND quantity > 0) as low_stock_total,
		(SELECT COUNT(*) FROM stocks WHERE quantity = 0) as no_stock_totalD
		`
	err := db.DB.QueryRow(widgetQuery).Scan(&DashboardData.ProductTotal, &DashboardData.StockTotal, &DashboardData.LowStockTotal, &DashboardData.NoStockTotal)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to dashboard data: "+err.Error())
		return
	}

	chartQuery := `SELECT d::date AS date, COALESCE(COUNT(tr.id), 0) AS stock_out_total
					FROM generate_series(
         				CURRENT_DATE - INTERVAL '6 days',
         				CURRENT_DATE,
         				INTERVAL '1 day'
     				) d
					LEFT JOIN transactions tr ON date_trunc('day', tr.created_at) = d::date AND tr.move_type = 'OUT'
					GROUP BY d
					ORDER BY d;`
	dbRows, err := db.DB.Query(chartQuery)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to dashboard data: "+err.Error())
		return
	}
	defer dbRows.Close()

	var chartData []map[string]interface{}
	for dbRows.Next() {
		var date string
		var stockOutTotal int
		err := dbRows.Scan(&date, &stockOutTotal)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to scan chart data: "+err.Error())
			return
		}
		chartData = append(chartData, map[string]interface{}{
			"date":  date,
			"total": stockOutTotal,
		})
	}

	DashboardData.ChartActivityData = chartData
	utils.RespondSuccess(w, DashboardData, "Dashboard data fetched successfully")

}
