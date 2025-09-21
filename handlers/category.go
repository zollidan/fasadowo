package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/fasadowo/models"
	"github.com/zollidan/fasadowo/utils"
	"gorm.io/gorm"
)


type CategoryHandler struct {
	DB *gorm.DB
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "categoryID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	category, err = gorm.G[models.Category](h.DB).Preload("Subcategories", nil).Where("id = ?", id).First(context.Background())
	if err == gorm.ErrRecordNotFound {
		utils.WriteError(w, http.StatusNotFound, "Category not found")
		return
	} 
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error fetchig category")
		return
	}
	
	json.NewEncoder(w).Encode(category)
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
