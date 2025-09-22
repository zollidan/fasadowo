package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/zollidan/fasadowo/config"
	"github.com/zollidan/fasadowo/database"
	"github.com/zollidan/fasadowo/handlers"
	"github.com/zollidan/fasadowo/httpmiddleware"
	"github.com/zollidan/fasadowo/models"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	cfg := config.New()

	// Инициализация JWT (секрет лучше хранить в .env или config)
	tokenAuth = jwtauth.New("HS256", []byte("devchik"), nil)

	db := database.InitDatabase()

	authHandler := handlers.AuthHandler{DB: db, TokenAuth: tokenAuth}
	categoryHandler := handlers.CategoryHandler{DB: db}
	subcategoriesHandler := handlers.SubcategoryHandler{DB: db}
	productHandler := handlers.ProductHandler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Get("/user/info", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				user, _ := gorm.G[models.User](db).Where("id = ?", claims["user_id"]).First(context.Background())
				json.NewEncoder(w).Encode(map[string]string{"email": user.Email})
			})
			r.With(httpmiddleware.RequireRole("admin")).Get("/admin", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				w.Write([]byte(fmt.Sprintf("welcome admin, user_id=%v", claims["user_id"])))
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.LoginUser)
			r.Post("/register", authHandler.RegisterUser)
		})

		r.Route("/category", func(r chi.Router) {
			r.Get("/", categoryHandler.ListCategory)
			r.Route("/{categoryID}", func(r chi.Router) {
				r.Get("/", categoryHandler.GetCategory)
			})
		})

		r.Route("/subcategory", func(r chi.Router) {
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

	log.Printf("Server is running on addr -> http://%s\n", cfg.ServerAddress())
	http.ListenAndServe(cfg.ServerAddress(), r)
}
