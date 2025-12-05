package main

import (
	"net/http"

	"github.com/Arrafll/StockLab-Go/internal/routes"
	_ "github.com/Arrafll/StockLab-Go/internal/routes"
)

// @title StockLab
// @version 1.0
// @description Dokumentasi API StockLab
// @contact.name Andre R
// @contact.email andrerafli83@gmail.com
// @license.name dreowsy
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" JWT
func main() {

	route := routes.RegisterRoutes()
	http.ListenAndServe(":8080", route)
}
