package main

import (
	"net/http"

	"github.com/Arrafll/setoko-go.git/internal/routes"
	_ "github.com/Arrafll/setoko-go.git/internal/routes"
)

func main() {
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	http.ListenAndServe(":8080", mux)
}
