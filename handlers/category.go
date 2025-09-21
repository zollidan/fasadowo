package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zollidan/fasadowo/models"
	"gorm.io/gorm"
)


type CategoryHandler struct {
	DB *gorm.DB
}
func (h *CategoryHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	if err := h.DB.Find(&categories).Error; err != nil {
		http.Error(w, "failed to fetch category", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "failed to serialize categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
