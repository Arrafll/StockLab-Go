package routes

import (
	"net/http"

	_ "github.com/Arrafll/StockLab-Go/docs" // <-- wajib ada
	"github.com/Arrafll/StockLab-Go/internal/config"
	authService "github.com/Arrafll/StockLab-Go/internal/services/auth"
	categoryService "github.com/Arrafll/StockLab-Go/internal/services/category"
	productService "github.com/Arrafll/StockLab-Go/internal/services/product"
	userService "github.com/Arrafll/StockLab-Go/internal/services/user"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(cfg *config.Config) http.Handler {

	r := chi.NewRouter()

	// Swagger UI route
	r.Get("/stocklab-api/documentation/*", httpSwagger.Handler(
		httpSwagger.URL("/stocklab-api/documentation/doc.json"), // URL relatif ke endpoint JSON
	))

	// API Version 1
	r.Route("/v1", func(r chi.Router) {
		// Login Routes
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			authService.Login(w, r, cfg)
		})
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Use(authService.JWTMiddleware(cfg)) // middleware JWT
			r.Get("/", userService.GetUserList)
			r.Post("/create", userService.CreateUser)
			r.Put("/update/{id}", userService.UpdateUser)
			r.Get("/detail/{id}", userService.GetUserDetail)
			r.Delete("/delete/{id}", userService.DeleteUser)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Use(authService.JWTMiddleware(cfg)) // middleware JWT
			r.Get("/", categoryService.GetCategoryList)
			r.Post("/create", categoryService.CreateCategory)
			r.Put("/update/{id}", categoryService.UpdateCategory)
			r.Delete("/delete/{id}", categoryService.DeleteCategory)
		})

		r.Route("/products", func(r chi.Router) {
			r.Use(authService.JWTMiddleware(cfg)) // middleware JWT
			r.Get("/", productService.GetProductList)
			r.Post("/create", productService.CreateProduct)
			r.Get("/detail/{id}", productService.GetProductDetail)
			r.Delete("/delete/{id}", productService.DeleteProduct)
			r.Patch("/update/{id}", productService.UpdateProduct)
		})

	})

	return r

}
