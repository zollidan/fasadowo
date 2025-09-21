package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/fasadowo/models"
	"github.com/zollidan/fasadowo/utils"
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

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "productID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := utils.GetByID[models.Product](w, h.DB, id)
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(product)
}