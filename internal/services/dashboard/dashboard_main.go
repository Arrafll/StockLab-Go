package services

import (
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/db"
	"github.com/Arrafll/StockLab-Go/internal/utils"
)

// Dashboard blueprint
type DashboardData struct {
	ProductTotal         int                      `json:"product_total"`
	StockTotal           int                      `json:"stock_total"`
	LowStockTotal        int                      `json:"low_stock"`
	NoStockTotal         int                      `json:"no_stock"`
	ChartActivityDataIn  []map[string]interface{} `json:"chart_activity_data_in"`
	ChartActivityDataOut []map[string]interface{} `json:"chart_activity_data_out"`
}

type DashboardSuccessResp struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    DashboardData `json:"data"`
}

type DashboardFailResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func DashboardMain(w http.ResponseWriter, r *http.Request) {
	// ✅ PAKAI camelCase
	var dashboardData DashboardData

	// Widget data
	widgetQuery := `
		SELECT 
			(SELECT COUNT(*) FROM products),
			(SELECT COALESCE(SUM(quantity),0) FROM stocks),
			(SELECT COUNT(*) FROM stocks WHERE quantity < 10 AND quantity > 0),
			(SELECT COUNT(*) FROM stocks WHERE quantity = 0)
	`

	err := db.DB.QueryRow(widgetQuery).Scan(
		&dashboardData.ProductTotal,
		&dashboardData.StockTotal,
		&dashboardData.LowStockTotal,
		&dashboardData.NoStockTotal,
	)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Chart IN
	dashboardData.ChartActivityDataIn, err = DashboardChartByMoveType("IN")
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Chart OUT
	dashboardData.ChartActivityDataOut, err = DashboardChartByMoveType("OUT")
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, dashboardData, "Dashboard data fetched successfully")
}

func DashboardChartByMoveType(moveType string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			d::date AS date,
			COALESCE(COUNT(tr.id), 0) AS total
		FROM generate_series(
			CURRENT_DATE - INTERVAL '6 days',
			CURRENT_DATE,
			INTERVAL '1 day'
		) d
		LEFT JOIN transactions tr
			ON date_trunc('day', tr.created_at) = d::date
			AND tr.move_type = $1
		GROUP BY d
		ORDER BY d;
	`

	rows, err := db.DB.Query(query, moveType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var date string
		var total int

		if err := rows.Scan(&date, &total); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"date":  date,
			"total": total,
		})
	}

	return result, nil
}
