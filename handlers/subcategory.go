package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zollidan/fasadowo/models"
	"gorm.io/gorm"
)

type SubcategoryHandler struct {
	DB *gorm.DB
}
func (h *SubcategoryHandler) ListSubcategory(w http.ResponseWriter, r *http.Request) {
	var subcategories []models.Subcategory
	if err := h.DB.Find(&subcategories).Error; err != nil {
		http.Error(w, "failed to fetch subcategory", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(subcategories)
	if err != nil {
		http.Error(w, "failed to serialize subcategories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
