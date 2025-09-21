package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zollidan/fasadowo/database"
	"github.com/zollidan/fasadowo/handlers"
)

func main() {
	db := database.InitDatabase()

	authHandler := handlers.AuthHandler{DB: db}
	categoryHandler := handlers.CategoryHandler{DB: db}
	subcategoriesHandler := handlers.SubcategoryHandler{DB: db}
	productHandler := handlers.ProductHandler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.RegisterUser)
		})
		r.Route("/category", func(r chi.Router) {
			r.Get("/", categoryHandler.ListCategory)
			r.Route("/{categoryID}", func(r chi.Router) {
				r.Get("/", categoryHandler.GetCategory)
			})
		})
		r.Route("/subcategory", func (r chi.Router)  {
			r.Get("/", subcategoriesHandler.ListSubcategory)
			r.Route("/{subcategoryID}", func(r chi.Router) {
				r.Get("/", subcategoriesHandler.GetSubcategory)
			})
		})
		r.Route("/product", func(r chi.Router) {
			r.Get("/", productHandler.ListProducts)
			r.Route("/{productID}", func(r chi.Router) {
				r.Get("/", productHandler.GetProduct)
			})
		})
	})

	http.ListenAndServe("localhost:3000", r)
}
