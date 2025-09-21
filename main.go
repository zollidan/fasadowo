package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zollidan/fasadowo/database"
	"github.com/zollidan/fasadowo/handlers"
	"github.com/zollidan/fasadowo/config"
)

func main() {
	db := database.InitDatabase()

	categoryHandler := handlers.CategoryHandler{DB: db}
	subcategoriesHandler := handlers.SubcategoryHandler{DB: db}
	productHandler := handlers.ProductHandler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/category", func(r chi.Router) {
			r.Get("/", categoryHandler.ListCategory)
		})
		r.Route("/subcategory", func (r chi.Router)  {
			r.Get("/", subcategoriesHandler.ListSubcategory)
		})
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.ListProducts)
		})
	})

	http.ListenAndServe("localhost:3000", r)
}
