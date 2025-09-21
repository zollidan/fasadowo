package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zollidan/fasadowo/models"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	if err := h.DB.Find(&products).Error; err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "failed to serialize products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
