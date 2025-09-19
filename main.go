package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	ctx context.Context
)

func listProducts(w http.ResponseWriter, r *http.Request) {
	product, err := gorm.G[Product](db).First(ctx)
	if err != nil {
		writeError()
	}

	jsonData, err := json.Marshal(product)
	if err != nil {
		writeError()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
}

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	ctx = context.Background()

	db.AutoMigrate(&Product{})

	err = gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100})
	if err != nil {
		log.Fatal("Error creating test product")
	}

	fmt.Println("Database initialized.")
}

func main() {

	initDatabase()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Route("/api", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Get("/", listProducts)                                
		})	                            
	})

	http.ListenAndServe("localhost:3000", r)
}